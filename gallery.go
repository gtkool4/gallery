// Package gallery creates a usable GTK4 widget gallery window.
package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gdkpixbuf/v2"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	externglib "github.com/diamondburned/gotk4/pkg/core/glib"

	"github.com/gtkool4/grun"
	"github.com/gtkool4/gtkelp/buildhelp"
	"github.com/gtkool4/gtkelp/gtknew"
)

const boxMargin = 2

const imageSource = "https://raw.githubusercontent.com/golang-samples/gopher-vector/master/"

var files = map[string][]byte{}

//
//--------------------------------------------------------------[ LAUNCH APP ]--

var gapp = &grun.App{
	ID:     "com.github.gtkool4.gallery",
	Title:  "GTK4 Gallery",
	Width:  800,
	Height: 800,
}

func main() {
	gapp.Run(func() gtk.Widgetter {
		// Preload pixbuf data for iconview.
		files["gotk4.png"] = downloadFile("https://avatars.githubusercontent.com/u/13782055?s=200&v=4")
		files["gopher-front.png"] = downloadFile(imageSource + "gopher-front.png")
		files["gopher-side.png"] = downloadFile(imageSource + "gopher-side_color.png")
		files["gopher.png"] = downloadFile(imageSource + "gopher.png")

		// Create widgets groups.
		widgets := make([]gtk.Widgetter, 6)
		titles := []string{"Displays", "Buttons", "Entries", "Containers", "Windows"}
		for i, list := range []Group{listDisplays, listButtons, listEntries, listContainers, listWindows} {
			widgets[i] = list.Widgets(titles[i])
		}
		widgets[5] = NewCustomWidgetStarted()

		box := gtknew.VBox(10, widgets...)
		return gtknew.ScrolledWindow(box)
	})
}

//
//------------------------------------------------------------[ WIDGETS LIST ]--

// Group defines a group of widget makers.
type Group []struct {
	Name string
	Make func() gtk.Widgetter
}

// Widgets creates a box with all widgets in the group.
func (l Group) Widgets(title string) gtk.Widgetter {
	var widgets []gtk.Widgetter
	for _, item := range l {
		widgets = append(widgets, gtknew.Frame(item.Name, item.Make()))
	}
	isWide := (title == "Containers")
	return gtknew.Frame(title, newContainer(isWide, widgets...))
}

var listDisplays = Group{
	{"Label", newLabel},
	{"Spinner", newSpinner},
	{"StatusBar", newStatusBar},
	{"LevelBar", newLevelBar},
	{"ProgressBar", newProgressBar},
	{"InfoBar", newInfoBar},
	{"ScrollBar", newScrollbar},
	{"Image", newImage},
	{"Picture", newPicture},
	{"Separator", newSeparator},
	{"TextView", newTextView},
	{"Scale", newScale},
	{"DrawingArea", newDrawingArea},
	{"Video", newVideo},
	{"MediaControls", newMediaControls},
	{"WindowControls", newWindowControls},
	{"MenuBar", newMenuBar},
	{"Calendar", newCalendar},
	{"EmojiChooser", placeholder},
	{"Menu", placeholder},
}

var listButtons = Group{
	{"Button", newButton},
	{"ToggleButton", newToggleGroup},
	{"LinkButton", newLinkButton},
	{"CheckButton", newCheckButton},
	{"RadioButton", newRadioGroup},
	{"MenuButton", newMenuButton},
	{"LockButton", newLockButton},
	{"VolumeButton", newVolumeButton},
	{"Switch", newSwitch},
	{"ComboBox", newComboBox},
	{"ComboBoxText", newComboBoxText},
	{"DropDown", newDropDown},
	{"ColorButton", newColorButton},
	{"FontButton", newFontButton},
	{"ApplicationButton", newApplicationButton},
}

var listEntries = Group{
	{"Entry", newEntry},
	{"SearchEntry", newSearchEntry},
	{"PasswordEntry", newPasswordEntry},
	{"Spinbutton", newSpinButton},
	{"EditableLabel", newEditableLabel},
}

var listContainers = Group{
	{"Box", newBox},
	{"Grid", newGrid},
	{"CenterBox", newCenterBox},
	{"ScrolledWindow", newScrolledWindow},
	{"Paned", newPaned},
	{"Frame", newFrame},
	{"Expander", newExpander},
	{"SearchBar", newSearchBar},
	{"ActionBar", newActionBar},
	{"HeaderBar", newHeaderBar},
	{"Notebook", newNotebook},
	{"ListBox", newListBox},
	{"FlowBox", newFlowBox},
	{"TreeView", newTreeView},
	{"Iconview", newIconview},
	{"Overlay", newOverlay},
	{"StackSwitcher", newStackSwitcher},
	{"StackSidebar", newStackSidebar},
	{"PopOver", newPopOver},
}

var listWindows = Group{
	{"Window", newWindow},
	{"Dialog", newDialog},
	{"CustomDialog", newCustomDialog},
	{"MessageDialog", newMessageDialog},
	{"AboutDialog", newAboutDialog},
	{"Assistant", newAssistant},
	{"PageSetupDialog", newPrintPageSetupDialog},
	{"PrintDialog", newPrintDialog},
	{"ShortcutsWindow", newShortcutsWindow},
	{"ColorChooser", newColorChooser},
	{"FileChooser", newFileChooser},
	{"FontChooser", newFontChooser},
	{"AppchooserDialog", newAppchooserDialog},
}

//
//---------------------------------------------------------[ DISPLAY WIDGETS ]--

func newLabel() gtk.Widgetter {
	return gtknew.VBox(boxMargin, gtk.NewLabel("Label"))
}

func newSpinner() gtk.Widgetter {
	w := gtk.NewSpinner()
	w.Start()
	return gtknew.VBox(boxMargin, w)
}

func newStatusBar() gtk.Widgetter {
	w := gtk.NewStatusbar()
	w.Push(0, "Statusbar")
	return w
}

func newLevelBar() gtk.Widgetter {
	w := gtk.NewLevelBar()
	w.SetValue(0.5)
	return w
}

func newProgressBar() gtk.Widgetter {
	w := gtk.NewProgressBar()
	w.SetFraction(0.5)
	return w
}

func newInfoBar() gtk.Widgetter {
	w := gtk.NewInfoBar()
	w.AddChild(gtk.NewLabel("InfoBar"))
	w.SetShowCloseButton(true)
	return gtknew.VBox(boxMargin, w)
}

func newScrollbar() gtk.Widgetter {
	a := gtk.NewAdjustment(1, 0, 200, 1, 10, 10)
	return gtk.NewScrollbar(gtk.OrientationHorizontal, a)
}

func newImage() gtk.Widgetter {
	w := gtk.NewImageFromIconName("face-cool")
	w.SetIconSize(gtk.IconSizeLarge)
	return w
}

func newPicture() gtk.Widgetter {
	pix := pixbufLoader(files["gotk4.png"])
	return gtk.NewPictureForPixbuf(pix)
}

func newSeparator() gtk.Widgetter {
	sh := gtk.NewSeparator(gtk.OrientationHorizontal)
	sv := gtk.NewSeparator(gtk.OrientationVertical)

	boxV := gtknew.VBox(20, // Vertical box for the vertical separator
		gtknew.HBox(0, sv), // empty box to have something before the separator
	)
	boxV.SetHExpand(true) // Expand full size so we can see our vertical separator

	return gtknew.HBox(40, // Horizontal box for both separators
		gtknew.HBox(0), // empty box to have something before the separator
		sh,
		boxV,
	)
}

func newTextView() gtk.Widgetter {
	tv := gtk.NewTextView()
	tv.SetVExpand(true)
	tv.SetHExpand(true)

	buffer := tv.Buffer()
	buffer.SetText("Text View\nis multiline", -1)

	return gtknew.VBox(boxMargin, tv)
}

func newScale() gtk.Widgetter {
	w := gtk.NewScaleWithRange(gtk.OrientationHorizontal, 0, 1, 0.1)
	w.SetValue(0.5)
	return gtknew.VBox(boxMargin, w)
}

func newDrawingArea() gtk.Widgetter {
	w := gtk.NewDrawingArea()
	w.SetContentWidth(100)
	w.SetContentHeight(100)
	w.SetDrawFunc(func(area *gtk.DrawingArea, cr *cairo.Context, width int, height int) {
		w, h := float64(width), float64(height)
		for i, color := range [][3]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}} {
			cr.Arc(w*(float64(i+1))/4, h/2, math.Min(w, h)/2, 0, 2*math.Pi)
			cr.SetSourceRGBA(color[0], color[1], color[2], 0.7)
			cr.Fill()
		}
	})
	return gtknew.VBox(boxMargin, w)
}

func newVideo() gtk.Widgetter {
	w := gtk.NewVideo()
	return w
}

func newMediaControls() gtk.Widgetter {
	w := gtk.NewMediaControls(gtk.NewMediaFile())
	return w
}

func newWindowControls() gtk.Widgetter {
	w := gtk.NewWindowControls(gtk.PackStart)
	w.SetDecorationLayout("icon:minimize,maximize,close")
	return w
}

func newMenuBar() gtk.Widgetter {
	return placeholder()
	// w:=gtk.NewMenuBar()
	// return w
}

func newCalendar() gtk.Widgetter {
	return gtk.NewCalendar()
}

// gtk.EmojiChooser.Realize is ambiguous
// cannot use w (type *gtk.EmojiChooser) as type gtk.Widgetter in argument to toBox:
// *gtk.EmojiChooser does not implement gtk.Widgetter (missing Realize method)
func newEmojiChooser() gtk.Widgetter {
	return placeholder()
	// w:=gtk.NewEmojiChooser()
	// return gtknew.VBox(boxMargin,w)
}

func newMenu() gtk.Widgetter { return placeholder() }

//
//---------------------------------------------------------[ BUTTONS WIDGETS ]--

func newButton() gtk.Widgetter {
	w := gtk.NewButtonWithLabel("Button")
	w.Connect("clicked", callPrint("button 1 clicked"))

	x := gtk.NewButtonFromIconName("preferences-system")
	x.Connect("clicked", callPrint("button 2 clicked"))
	return gtknew.VBox(boxMargin, w, x)
}

func newCheckButton() gtk.Widgetter {
	btn1 := gtk.NewCheckButtonWithLabel("CheckButton")
	btn2 := gtk.NewCheckButtonWithLabel("Not checked")
	btn1.SetActive(true)
	btn1.Connect("toggled", callPrint("check 1 toggled"))
	btn2.Connect("toggled", callPrint("check 2 toggled"))
	return gtknew.VBox(boxMargin, btn1, btn2)
}

func newLinkButton() gtk.Widgetter {
	w := gtk.NewLinkButtonWithLabel("https://golang.org/", "Link Button")
	w.Connect("clicked", callPrint("link button clicked"))
	return gtknew.VBox(boxMargin, w)
}

func newToggleGroup() gtk.Widgetter {
	var group *gtk.ToggleButton // Reference to previous button, so we can add the new one in the same group.
	box := gtknew.VBox(boxMargin)
	for i, txt := range []string{"Toggle", "Button"} {
		btn := gtk.NewToggleButtonWithLabel(txt)
		if group == nil {
			btn.SetActive(true)
		} else {
			btn.SetGroup(group)
		}

		group = btn
		box.Append(btn)

		i := i // We're in a loop, so we need to make a static copy of the index for the callback.
		btn.Connect("toggled", func() { fmt.Println("toggle button", i, btn.Active()) })
	}
	return box
}

func newRadioGroup() gtk.Widgetter {
	var group *gtk.CheckButton // Reference to previous button, so we can add the new one in the same group.
	box := gtknew.VBox(boxMargin)
	for i, txt := range []string{"Radio Button", "second option"} {
		btn := gtk.NewCheckButtonWithMnemonic(txt)
		if group == nil {
			btn.SetActive(true) // Set the first button in clicked state.
		} else {
			btn.SetGroup(group) // Add other buttons in a group with the first one.
		}

		group = btn
		box.Append(btn)

		i := i // We're in a loop, so we need to make a static copy of the index for the callback.
		btn.Connect("toggled", func() { fmt.Println("radio button", i, btn.Active()) })
	}
	return box
}

func newMenuButton() gtk.Widgetter {
	btn := gtk.NewMenuButton()
	btn.Connect("activate", callPrint("menu button value changed")) // since gtk 4.4

	isMaximized := gapp.Win.IsMaximized()
	vMax := glib.NewVariantBoolean(isMaximized)
	actFullScreen := gio.NewSimpleActionStateful("fullscreen", nil, vMax)
	actQuit := gio.NewSimpleAction("quit", nil)

	actFullScreen.Connect("change-state", func() { // Args: *gio.SimpleAction, *glib.Variant  (the variant crash ATM)
		fmt.Println("menu actFullScreen", vMax.Boolean())
		newval := !gapp.Win.IsMaximized()
		if newval {
			gapp.Win.Maximize()
		} else {
			gapp.Win.Unmaximize()
		}
		vMax = glib.NewVariantBoolean(newval)
		actFullScreen.SetState(vMax)
	})

	actQuit.Connect("activate", func() {
		fmt.Println("menu actQuit")
		gapp.App.Quit()
	})

	gapp.Win.AddAction(actFullScreen)
	gapp.App.AddAction(actQuit)

	menu := gio.NewMenu()
	menu.Append("FullScreen", "win.fullscreen") //"app.action")
	menu.Append("Quit", "app.quit")
	btn.SetDirection(gtk.ArrowNone) // Hide the button arrow and restore the default button icon.
	btn.SetMenuModel(menu)

	return gtknew.VBox(boxMargin, &btn.Widget)
}

func newLockButton() gtk.Widgetter {
	w := gtk.NewLockButton(gio.NewSimplePermission(false))
	w.Connect("clicked", callPrint("lock button clicked"))
	return w
}

func newVolumeButton() gtk.Widgetter {
	w := gtk.NewVolumeButton()
	w.Connect("value-changed", callPrint("volume button value changed"))
	return gtknew.VBox(boxMargin, w)
}

func newSwitch() gtk.Widgetter {
	btn1 := gtk.NewSwitch()
	btn2 := gtk.NewSwitch()
	btn1.SetActive(true)
	btn1.Connect("activate", callPrint("switch 1 activated"))
	btn2.Connect("activate", callPrint("switch 2 activated"))

	box := gtk.NewCenterBox()
	box.SetCenterWidget(gtknew.VBox(boxMargin, btn1, btn2))
	return box
}

// ModelCB defines data fields order for the ListStore (TreeView, IconView, ComboBox).
const ( // Must match the ListStore declaration order.
	ModelCBRef = iota
	ModelCBText
	ModelCBTooltip
	ModelCBIcon
)

func newComboBox() gtk.Widgetter {
	model := gtk.NewListStore([]externglib.Type{
		externglib.TypeString,
		externglib.TypeString,
	})

	for _, item := range [][]string{
		{"0", "ComboBox"},
		{"1", "with"},
		{"2", "choices"},
	} {
		insertWithValues(model, map[int]interface{}{
			ModelCBRef:  item[0],
			ModelCBText: item[1],
		})
	}

	list := gtk.NewComboBoxWithModel(model)
	list.SetIDColumn(ModelCBRef)
	list.SetEntryTextColumn(ModelCBText)

	// Don't forget to add a cell renderer.
	cellText := gtk.NewCellRendererText()
	list.PackStart(cellText, false)
	list.AddAttribute(cellText, "markup", ModelCBText)
	list.SetActiveID("0")

	// Get values examples.
	iter, ok := list.ActiveIter()
	if ok {
		gval := model.Value(&iter, ModelCBText)
		if gval.GoValue() != "ComboBox" {
			fmt.Println("error combobox iter value is not", "ComboBox")
		}
	}

	content := []string{}
	model.Foreach(func(model gtk.TreeModeller, path *gtk.TreePath, iter *gtk.TreeIter) (stop bool) {
		gval := model.Value(iter, ModelCBText)
		content = append(content, gval.GoValue().(string))
		return false
	})
	if strings.Join(content, " ") != "ComboBox with choices" {
		fmt.Println("error combobox values modified")
	}

	return list
}

func newComboBoxText() gtk.Widgetter {
	w := gtk.NewComboBoxText()
	w.AppendText("ComboBoxText")
	w.AppendText("is easier")
	w.AppendText("to implement")
	w.SetActive(0)
	w.Connect("changed", callPrint("combo box text selection changed"))
	return gtknew.VBox(boxMargin, w)
}

func newDropDown() gtk.Widgetter {
	return gtknew.VBox(boxMargin, gtk.NewDropDownFromStrings([]string{"Drop Down", "is better", "than Combo Box"}))
}

func newColorButton() gtk.Widgetter {
	color := gdk.NewRGBA(1, 0, 0, 1)
	w := gtk.NewColorButton()
	w.SetRGBA(&color)
	w.Connect("color-set", callPrint("color button color set"))
	return gtknew.VBox(boxMargin, w)
}

func newFontButton() gtk.Widgetter {
	w := gtk.NewFontButton()
	w.Connect("activate", callPrint("font button activated")) // since gtk 4.4
	w.Connect("font-set", func() { fmt.Println("font button new font set:", w.FontDesc()) })
	return gtknew.VBox(boxMargin, &w.Widget)
}

func newApplicationButton() gtk.Widgetter {
	w := gtk.NewAppChooserButton("video/avi")
	w.Connect("changed", callPrint("application button changed"))
	return gtknew.VBox(boxMargin, w)
}

//
//-----------------------------------------------------------------[ ENTRIES ]--

func newEntry() gtk.Widgetter {
	w := gtk.NewEntry()
	w.Buffer().SetText("Entry", -1)
	w.Connect("changed", func() { fmt.Printf("entry changed: '%s'\n", w.Buffer().Text()) })
	return gtknew.VBox(boxMargin, w)
}

func newSearchEntry() gtk.Widgetter {
	w := gtk.NewSearchEntry()
	w.SetText("SearchEntry")
	w.Connect("search-changed", func() { fmt.Printf("search entry changed: '%s'\n", w.Text()) })
	return gtknew.VBox(boxMargin, w)
}

func newPasswordEntry() gtk.Widgetter {
	w := gtk.NewPasswordEntry()
	w.SetText("PasswordEntry")
	w.SetShowPeekIcon(true)
	w.Connect("changed", func() { fmt.Printf("password entry changed: '%s'\n", w.Text()) })
	return gtknew.VBox(boxMargin, w)
}

func newSpinButton() gtk.Widgetter {
	w := gtk.NewSpinButtonWithRange(1, 100, 1)
	w.SetValue(42)
	w.Connect("changed", func() { fmt.Printf("spin button changed: '%s'\n", w.Text()) })
	return gtknew.VBox(boxMargin, w)
}

func newEditableLabel() gtk.Widgetter {
	w := gtk.NewEditableLabel("EditableLabel")
	w.StartEditing()
	w.Connect("changed", func() { fmt.Printf("editable entry changed: '%s'\n", w.Text()) })
	return gtknew.VBox(boxMargin, w)
}

//
//--------------------------------------------------------------[ CONTAINERS ]--

func newBox() gtk.Widgetter {
	h := gtk.NewBox(gtk.OrientationHorizontal, boxMargin)
	h.Append(gtknew.Frame("Horizontal"))
	h.Append(gtknew.Frame("-----"))
	h.Append(gtknew.Frame("-----"))

	v := gtk.NewBox(gtk.OrientationVertical, boxMargin)
	v.Append(gtknew.Frame("Vertical"))
	v.Append(gtknew.Frame("-----"))
	v.Append(gtknew.Frame("-----"))

	h.SetHExpand(true)
	v.SetHExpand(true)
	h.SetHAlign(gtk.AlignCenter)
	v.SetHAlign(gtk.AlignCenter)
	return newHBoxExpand(h, v)
}

func newGrid() gtk.Widgetter {
	w := gtk.NewGrid()
	for i := 0; i < 6; i++ {
		w.Attach(gtknew.Frame(strconv.Itoa(int(i)+1)), i%3, i/3, 1, 1)
	}
	return w
}

func newCenterBox() gtk.Widgetter {
	w := gtk.NewCenterBox()
	w.SetStartWidget(gtknew.Frame("Left"))
	w.SetCenterWidget(gtknew.Frame("Center"))
	w.SetEndWidget(gtknew.Frame("Right"))
	return w
}

func newScrolledWindow() gtk.Widgetter {
	box := gtknew.VBox(400)
	box.Append(gtk.NewLabel("ScrolledWindow"))
	box.Append(gtk.NewLabel("Bottom"))
	w := gtk.NewScrolledWindow()
	w.SetChild(gtknew.VBox(boxMargin, box))
	w.SetHasFrame(true)
	return w
}

func newPaned() gtk.Widgetter {
	w := gtk.NewPaned(gtk.OrientationHorizontal)
	w.SetStartChild(gtknew.Frame("Left"))
	w.SetEndChild(gtknew.Frame("Right"))
	return w
}

func newFrame() gtk.Widgetter {
	w := gtk.NewFrame("Frame")
	w.SetChild(gtk.NewLabel("with child"))
	return gtknew.VBox(boxMargin, w)
}

func newExpander() gtk.Widgetter {
	w := gtk.NewExpander("Expander")
	w.SetChild(gtk.NewLabel("This was hidden"))
	return w
}

func newSearchBar() gtk.Widgetter {
	e := gtk.NewSearchEntry()
	e.SetHExpand(true)
	e.SetText("SearchBar")
	w := gtk.NewSearchBar()
	w.SetChild(gtknew.VBox(boxMargin, e))
	w.SetShowCloseButton(true)
	w.SetSearchMode(true)
	return w
}

func newActionBar() gtk.Widgetter {
	cut := gtk.NewButtonFromIconName("edit-cut")
	copy := gtk.NewButtonFromIconName("edit-copy")
	paste := gtk.NewButtonFromIconName("edit-paste")
	close := gtk.NewButtonFromIconName("window-close")
	cut.SetTooltipText(("Cut"))
	copy.SetTooltipText(("Copy"))
	paste.SetTooltipText(("Paste"))
	close.SetTooltipText(("Close"))

	w := gtk.NewActionBar()
	w.PackStart(newHBoxExpand(cut, copy, paste))
	w.PackEnd(newHBoxExpand(close))
	w.SetHExpand(true)

	btn := gtk.NewToggleButtonWithLabel("Show")
	btn.SetTooltipText("Show or hide the action bar")
	btn.Connect("toggled", func() { w.SetRevealed(btn.Active()) })
	btn.SetActive(true)
	return newHBoxExpand(btn, w)
}

func newHeaderBar() gtk.Widgetter {
	w := gtk.NewHeaderBar()
	w.SetShowTitleButtons(true)
	w.SetDecorationLayout("icon:close")
	label := gtk.NewLabel("<b>HeaderBar</b>")
	label.SetUseMarkup(true)
	w.SetTitleWidget(label)
	return w
}

func newNotebook() gtk.Widgetter {
	w := gtk.NewNotebook()
	w.AppendPage(gtk.NewLabel("Notebook"), gtk.NewLabel("Page 1"))
	w.AppendPage(gtk.NewLabel("another"), gtk.NewLabel("Page 2"))
	w.AppendPage(gtk.NewLabel("page"), gtk.NewLabel("Page 3"))
	w.Connect("switch-page", callPrint("notebook page changed"))
	return w
}

func newListBox() gtk.Widgetter {
	b1 := newHBoxExpand()
	b2 := newHBoxExpand()
	b3 := newHBoxExpand()
	b1.SetHExpand(true)
	b2.SetHExpand(true)
	b3.SetHExpand(true)

	w := gtk.NewListBox()
	w.Append(newHBoxExpand(gtk.NewLabel("Line One"), b1, gtk.NewCheckButton()))
	w.Append(newHBoxExpand(gtk.NewLabel("Line Two"), b2, gtk.NewButtonWithLabel("2")))
	w.Append(newHBoxExpand(gtk.NewLabel("Line Three"), b3, gtk.NewEntry()))
	return w
}

func newFlowBox() gtk.Widgetter {
	w := gtk.NewFlowBox()
	w.Insert(gtk.NewLabel("Child One"), 0)
	w.Insert(gtk.NewButtonWithLabel("Child Two"), 1)
	w.Insert(gtk.NewCheckButtonWithLabel("Child Three"), 2)
	return w
}

func newTreeView() gtk.Widgetter {
	// Order in the model must match the const declaration order (ModelCB...)
	model := gtk.NewListStore([]externglib.Type{
		externglib.TypeString, //  REF
		externglib.TypeString, //  TEXT
		externglib.TypeString, //  COMMENT
		externglib.TypeInt,    //  VALUE
	})

	for _, data := range []struct { // Fill the model with values.
		ref, text, comment string
		value              float64
	}{
		{"0", "<b><big>TreeView</big></b>", "text 1", 20},
		{"1", "with", "text 2", 50},
		{"2", "texts", "<s>text 3</s>", 80},
	} {
		insertWithValues(model, map[int]interface{}{
			ModelCBRef:         data.ref,
			ModelCBText:        data.text,
			ModelCBTooltip:     data.comment,
			ModelCBTooltip + 1: data.value,
		})
	}

	// Create TreeView
	w := gtk.NewTreeViewWithModel(model)

	// Add editable text column
	cellText := gtk.NewCellRendererText()
	cellText.SetObjectProperty("editable", true)
	cellText.Connect("edited",
		func(_ *gtk.CellRendererText, path, text string) { fmt.Println(text) })
	columnText := gtk.NewTreeViewColumn()
	columnText.SetTitle("Name")
	columnText.SetResizable(true)

	columnText.PackEnd(cellText, false)
	columnText.AddAttribute(cellText, "markup", ModelCBText)
	w.AppendColumn(columnText)

	// Add simple text column
	cellTooltip := gtk.NewCellRendererText()
	columnTooltip := gtk.NewTreeViewColumn()
	columnTooltip.SetTitle("Comment")

	columnTooltip.PackEnd(cellTooltip, false)
	columnTooltip.AddAttribute(cellTooltip, "markup", ModelCBTooltip)
	w.AppendColumn(columnTooltip)

	// Add progress bar column
	cellProgress := gtk.NewCellRendererProgress()
	cellProgress.SetObjectProperty("value", 100)
	columnProgress := gtk.NewTreeViewColumn()

	columnProgress.PackEnd(cellProgress, false)
	columnProgress.AddAttribute(cellProgress, "value", ModelCBTooltip+1)
	w.AppendColumn(columnProgress)

	return &w.Widget
}

func newIconview() gtk.Widgetter {
	// Order in the model must match the const declaration order (ModelCB...)
	model := gtk.NewListStore([]externglib.Type{
		externglib.TypeInt,                   // REF. The key column can also be an int.
		externglib.TypeString,                // TEXT
		externglib.TypeString,                // TOLTIP
		externglib.TypeFromName("GdkPixbuf"), // ICON
	})

	for _, data := range []struct { // Fill the model with values.
		ref                 int
		text, tooltip, icon string
	}{
		{0, "<big>Iconview</big>", "tooltip 1", "gotk4.png"},
		{1, "with", "tooltip 2", "gopher-front.png"},
		{2, "icons", "tooltip 3", "gopher-side.png"},
		{3, "and <b>tooltips</b>", "tooltip 4", "gopher.png"},
	} {
		iter := model.Append()
		model.SetValue(&iter, ModelCBRef, externglib.NewValue(data.ref))
		model.SetValue(&iter, ModelCBText, externglib.NewValue(data.text))
		model.SetValue(&iter, ModelCBTooltip, externglib.NewValue(data.tooltip))

		content := files[data.icon]

		// TODO: find why this can't be loaded here (http.Transport panic).
		// Should be:
		// content := downloadFile(imageSource + data.icon)

		externglib.IdleAdd(func() {

			pixbuf := pixbufLoader(content)
			if pixbuf != nil {
				model.SetValue(&iter, ModelCBIcon, externglib.NewValue(pixbuf))
			}
		})
	}

	w := gtk.NewIconViewWithModel(model)
	w.SetActivateOnSingleClick(true)
	w.SetMarkupColumn(ModelCBText)
	w.SetTooltipColumn(ModelCBTooltip)
	w.SetPixbufColumn(ModelCBIcon)
	w.Connect("selection-changed", callPrint("iconview selection changed"))

	return gtknew.VBox(boxMargin, &w.Widget)
}

func newOverlay() gtk.Widgetter {
	w := gtk.NewOverlay()
	w.SetChild(gtk.NewLabel("Overlay"))
	w.AddOverlay(gtk.NewLabel("\non top"))
	return w
}

func newStack() *gtk.Stack {
	stack := gtk.NewStack()
	t1 := gtk.NewTextView()
	t2 := gtk.NewTextView()
	t3 := gtk.NewTextView()
	t1.Buffer().SetText("Page 1\n\n\n\n\ncontent", -1)
	t2.Buffer().SetText("Page 2", -1)
	t3.Buffer().SetText("Page 3", -1)
	stack.AddTitled(t1, "page1", "Page 1")
	stack.AddTitled(t2, "page2", "Page 2")
	stack.AddTitled(t3, "page3", "Page 3")
	return stack
}

func newStackSwitcher() gtk.Widgetter {
	stack := newStack()
	stack.SetTransitionType(gtk.StackTransitionTypeOverLeft)
	sw := gtk.NewStackSwitcher()
	sw.SetStack(stack)
	return gtknew.VBox(boxMargin, sw, stack)
}

func newStackSidebar() gtk.Widgetter {
	stack := newStack()
	stack.SetTransitionType(gtk.StackTransitionTypeOverUp)
	stack.SetVisibleChildName("page2")
	stack.SetHExpand(true)
	stack.SetVExpand(true)
	sw := gtk.NewStackSidebar()
	sw.SetStack(stack)

	box := gtknew.HBox(boxMargin)
	box.Append(sw)
	box.Append(stack)

	return box
}

func newPopOver() gtk.Widgetter { return placeholder() }

//
//-----------------------------------------------------------------[ WINDOWS ]--

func newWindow() gtk.Widgetter {
	return buttonAction("Window", "document-new", func() {
		win := gtk.NewApplicationWindow(gapp.App)
		// win := gtk.NewWindow() // Another option.
		win.SetChild(gtk.NewLabel("Hello gotk4"))
		win.SetDefaultSize(300, 200)
		win.Show()
	})
}

func newDialog() gtk.Widgetter {
	return buttonAction("Dialog", "document-open", func() {
		w := gtk.NewDialog()
		w.AddButton("_OK", int(gtk.ResponseAccept))
		w.AddActionWidget(gtk.NewButtonFromIconName("document-open"), 1)
		w.AddActionWidget(gtk.NewButtonFromIconName("document-save"), 2)
		w.AddButton("_Cancel", int(gtk.ResponseCancel))
		w.SetDefaultResponse(0)
		w.ContentArea().Append(gtk.NewLabel("Dialog\n\nwith buttons"))
		w.Connect("response", func(d *gtk.Dialog, resp int) { fmt.Println("dialog response", resp); w.Destroy() })

		w.Show()
	})
}

func newCustomDialog() gtk.Widgetter {
	return buttonAction("Custom Dialog", "document-properties", func() {
		b := buildhelp.NewFromString( // example copied from dialog documentation. TODO: improve
			`
<?xml version="1.0" encoding="UTF-8"?>
<interface>
   <object class="GtkDialog" id="dialog1">
     <child type="action">
       <object class="GtkButton" id="button_cancel"/>
     </child>
     <child type="action">
       <object class="GtkButton" id="button_ok">
       </object>
     </child>
     <action-widgets>
       <action-widget response="cancel">button_cancel</action-widget>
       <action-widget response="ok" default="true">button_ok</action-widget>
     </action-widgets>
   </object>
</interface>`)

		w := b.Dialog("dialog1")
		testError(b.Errors())
		w.Connect("response", func(d *gtk.Dialog, resp int) { fmt.Println("custom dialog response", resp); w.Destroy() })
		w.Show()
	})
}

func newMessageDialog() gtk.Widgetter { return placeholder() }

func newAboutDialog() gtk.Widgetter {
	return buttonAction("AboutDialog", "help-about", func() {
		w := gtk.NewAboutDialog()
		w.SetArtists([]string{"artists", "and", "others"})
		w.SetAuthors([]string{"authors", "and", "others"})
		w.SetComments("comments")
		w.SetCopyright("copyright")
		w.SetDocumenters([]string{"documenters", "and", "others"})
		// w.SetLicense("license") // Overridden by SetLicenseType
		// w.SetWrapLicense(true)
		w.SetLicenseType(gtk.LicenseMITX11)
		// w.SetLogo(gtk.NewImageFromIconName("document-new").Paintable()) // TODO: bug
		w.SetLogoIconName("document-new")
		w.SetProgramName("name")
		w.SetSystemInformation("system Information")
		w.SetTranslatorCredits("translator Credits")
		w.SetVersion("version")
		w.SetWebsite("website")
		w.SetWebsiteLabel("website Label")
		w.Show()
	})
}

func newAssistant() gtk.Widgetter {
	return buttonAction("Assistant", "system-help", func() {
		w := gtk.NewAssistant()
		w.SetTitle("Assistant")
		w.SetSizeRequest(400, 300)
		w.Connect("close", w.Destroy)
		w.Connect("cancel", w.Destroy)

		appendPage := func(complete bool, typ gtk.AssistantPageType, title, content string) {
			box := gtknew.VBox(0)
			box.Append(gtk.NewLabel(content))
			w.AppendPage(box)
			w.SetPageType(box, typ)
			w.SetPageTitle(box, title)
			if complete {
				w.SetPageComplete(box, complete)
			}
		}

		appendPage(true, gtk.AssistantPageIntro, "Introduction", "This is just a basic example\nof gtk Assistant.")
		appendPage(true, gtk.AssistantPageProgress, "Step 1 of 2", "This is Step 1.")
		appendPage(true, gtk.AssistantPageProgress, "Step 2 of 2", "This is Step 2.")
		appendPage(false, gtk.AssistantPageSummary, "Conclusion", "Conclusion.")
		w.Show()
	})
}

func newPrintPageSetupDialog() gtk.Widgetter {
	return buttonAction("PageSetupDialog", "document-page-setup", func() {
		gtk.PrintRunPageSetupDialogAsync(
			gtk.NewWindow(),
			gtk.NewPageSetup(),
			gtk.NewPrintSettings(),
			func(pageSetup *gtk.PageSetup) { fmt.Println("print page setup dialog closed") },
		)
	})
}

func newPrintDialog() gtk.Widgetter {
	// return buttonAction("PrintDialog", "document-print", func() {}) // "printer"
	return placeholder()
}

func newShortcutsWindow() gtk.Widgetter {
	return buttonAction("ShortcutsWindow", "preferences-desktop-keyboard", func() {
		b := buildhelp.NewFromFile("shortcuts-clocks.ui")
		w := b.ShortcutsWindow("shortcuts-clocks")
		testError(b.Errors())
		w.Show()
	})
}

func newColorChooser() gtk.Widgetter {
	w := gtk.NewColorChooserWidget()
	w.Connect("color-activated", func() { col := w.RGBA(); fmt.Println("color chooser new color activated:", col.String()) })
	return gtknew.Expander("ColorChooser", w)
}

func newFileChooser() gtk.Widgetter {
	w := gtk.NewFileChooserWidget(gtk.FileChooserActionOpen)
	w.SetSelectMultiple(true)
	w.Connect("up-folder", func() { fmt.Println("file chooser up-folder") })
	w.Connect("down-folder", func() { fmt.Println("file chooser down-folder") })
	return gtknew.Expander("FileChooser", w)
}

func newFontChooser() gtk.Widgetter {
	w := gtk.NewFontChooserWidget()
	w.Connect("font-activated", func() { fmt.Println("font chooser new font activated:", w.FontDesc()) })
	return gtknew.Expander("FontChooser", &w.Widget)
}

func newAppchooserDialog() gtk.Widgetter {
	w := gtk.NewAppChooserWidget("video/avi")
	w.Connect("application-selected", func() { fmt.Println("app chooser new application selected:", w.AppInfo().Name()) })
	return gtknew.Expander("AppchooserDialog", &w.Widget)
}

//
//-----------------------------------------------------------------[ CUSTOM WIDGET ]--

// CustomWidget provides an example of custom widget that behaves like an
// extended gtk.Box.
//
// It displays a timer and provides a switch that updates an icon and a text.
type CustomWidget struct {
	gtk.Box    // Extends the main container.
	sw         *gtk.Switch
	img        *gtk.Image
	labelState *gtk.Label
	labelTime  *gtk.Label
}

// NewCustomWidgetStarted creates and starts the custom widget.
func NewCustomWidgetStarted() *CustomWidget {
	w := NewCustomWidget()
	go w.Loop()
	return w
}

// NewCustomWidget creates an example widget with timer and switch.
func NewCustomWidget() *CustomWidget {
	box := &CustomWidget{
		Box:        *gtknew.HBox(10),
		sw:         gtk.NewSwitch(),
		img:        gtk.NewImage(),
		labelState: gtk.NewLabel("Custom Widget"),
		labelTime:  gtk.NewLabel(""),
	}

	box.labelTime.SetHAlign(gtk.AlignEnd)
	box.labelTime.SetHExpand(true)

	box.sw.Connect("activate", func() { fmt.Println("custom activate", box.Active()) }) // signal "activate" doesn't seem to work
	box.sw.Connect("state-set", box.switchToggled)
	box.SetActive(true)

	// Packing
	box.Append(gtk.NewLabel("Custom Widget:"))
	box.Append(box.sw)
	box.Append(box.img)
	box.Append(box.labelState)
	box.Append(box.labelTime)
	return box
}

// Widget Public API.

// Active returns the switch value.
func (w *CustomWidget) Active() bool { return w.sw.Active() }

// SetActive sets the switch value.
func (w *CustomWidget) SetActive(active bool) { w.sw.SetActive(active) }

// Loop starts the timer loop.
func (w *CustomWidget) Loop() {
	// When called from a go routine, use IdleAdd to run your gtk actions within
	// the gtk main loop.
	for t := range time.Tick(time.Second) {
		text := fmt.Sprintf("Last updated at %s.", t.Format(time.StampMilli))
		gtknew.Idle(func() { w.labelTime.SetLabel(text) })
	}
}

// Widget Private Callbacks.
func (w *CustomWidget) switchToggled() {
	// No need for IdleAdd as this is called from a gtk callback, when the switch
	// is toggled, but also at widget creation during the window creation call
	// with SetActive(true).
	text := map[bool]string{false: "Active", true: "Inactive"}
	icon := map[bool]string{false: "go-up", true: "go-down"}
	newValue := !w.Active() // reverse value as this is called before the change.
	w.img.SetFromIconName(icon[newValue])
	w.labelState.SetLabel(text[newValue])
	fmt.Println("custom widget switched to", newValue)
}

//
//-----------------------------------------------------------------[ COMMON ]--

func newContainer(isWide bool, items ...gtk.Widgetter) gtk.Widgetter {
	box := gtk.NewFlowBox()
	box.SetSelectionMode(gtk.SelectionNone)
	for i, item := range items {
		box.Insert(item, int(i))
	}
	if isWide {
		box.SetMaxChildrenPerLine(2)
	}
	return box
}

func newHBoxExpand(list ...gtk.Widgetter) gtk.Widgetter {
	box := gtknew.HBox(boxMargin, list...)
	box.SetHExpand(true)
	return box
}

func buttonAction(label, icon string, call func()) gtk.Widgetter {
	w := gtk.NewButton()
	w.Connect("clicked", call)
	img := gtk.NewImageFromIconName(icon)
	img.SetIconSize(gtk.IconSizeLarge)
	w.SetChild(gtknew.VBox(boxMargin, img, gtk.NewLabel(label)))
	return w
}

func insertWithValues(model *gtk.ListStore, data map[int]interface{}) {
	var keys []int
	var values []externglib.Value
	for k, v := range data {
		keys = append(keys, k)
		values = append(values, *externglib.NewValue(v))
	}

	model.InsertWithValues(-1, keys, values)
}

func downloadFile(url string) []byte {
	resp, e := http.Get(url)
	if e != nil {
		fmt.Printf("can't find icon url (%s): %s", url, e)
		return nil
	}
	byts, e := io.ReadAll(resp.Body)
	resp.Body.Close()
	if e != nil {
		fmt.Printf("can't download icon (%s): %s", url, e)
		return nil
	}
	return byts
}

func pixbufLoader(byts []byte) *gdkpixbuf.Pixbuf {
	load := gdkpixbuf.NewPixbufLoader()
	load.Connect("size-prepared", func(m *gdkpixbuf.PixbufLoader, w, h int) {
		m.SetSize(48, 48)
	})
	load.Write(byts)
	pix := load.Pixbuf()
	e := load.Close()
	if e != nil {
		fmt.Printf("can't close pixbuf loader: %s", e)
	}
	return pix
}

func placeholder() gtk.Widgetter { return gtk.NewLabel("TODO") }

func callPrint(args ...interface{}) func() { return func() { fmt.Println(args...) } }

func testError(errs grun.Errors) {
	if errs.IsError() {
		fmt.Println(errs.Error())
		panic("builder errors")
	}
}
