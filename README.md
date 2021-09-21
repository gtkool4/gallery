# GTK4 Widget Gallery example

This program provides a view similar to [The GTK-4.0 Widget Gallery](https://docs.gtk.org/gtk4/visual_index.html) page.
* It's also a running preview to test and choose your widgets.
* And code examples on how to use them.

## Displays
![GitHub Logo](https://raw.githubusercontent.com/gtkool4/assets/master/widgetimg/gallery-displays-20210919-1.png)
## Buttons
![GitHub Logo](https://raw.githubusercontent.com/gtkool4/assets/master/widgetimg/gallery-buttons-20210919-1.png)
## Entries
![GitHub Logo](https://raw.githubusercontent.com/gtkool4/assets/master/widgetimg/gallery-entries-20210919-1.png)
## Containers
![GitHub Logo](https://raw.githubusercontent.com/gtkool4/assets/master/widgetimg/gallery-containers-20210919-1.png)
## Windows
![GitHub Logo](https://raw.githubusercontent.com/gtkool4/assets/master/widgetimg/gallery-windows1-20210919-1.png)
![GitHub Logo](https://raw.githubusercontent.com/gtkool4/assets/master/widgetimg/gallery-windows2-20210919-1.png)
![GitHub Logo](https://raw.githubusercontent.com/gtkool4/assets/master/widgetimg/gallery-windows3-20210919-1.png)


## Missing widgets

* WindowControls : don't show
* EmojiChooser
* Menu
* PopOver
* MessageDialog
* PrintDialog (os dependent?)
* GLArea: Won't do here, too many deps. Will have a dedicated example.
 

## Problems

* SearchBar
  * Cannot use (*SearchBar).ConnectEntry(SearchEntry)
    * `cannot use variable of type *gtk.SearchEntry as gtk.Editabler: missing method Editable`
    * [example on the gnome repo](https:itlab.gnome.org/GNOME/gtk/-/blob/master/examples/search-bar.c)
* AboutDialog
  * Panics when trying to SetLogo(Paintable)
* Switch
  * Changing Switch.Connect("state-set") to ConnectAfter breaks the callback (in CustomWidget).
* Dialog
  * `Gtk-Message: GtkDialog mapped without a transient parent. This is discouraged.`
    * w.SetTransientFor(gapp.Win) : `cannot use gapp.Win (variable of type *gtk.ApplicationWindow) as *gtk.Window`
    * w.SetParent(gapp.Win) : `cannot use gapp.Win (variable of type *gtk.ApplicationWindow) as gtk.Widgetter`
* LockButton
  * can't unlock, maybe need a better gio.Permissioner (only found one usable)
* PixbufLoader
  * Would be nice to change the returns to be able to use as io.Writer (wrong type for method Write)
    * have func([]byte) error
    * want func([]byte) (n int, err error)
* MenuButton
  * missing callback args
  * using .Widget to prevent the naming conflict;
* FontButton
  * using .Widget to prevent the naming conflict
* TreeView
  * Using .Widget to prevent the naming conflict
  * When editing a cell: Gtk-CRITICAL :
     * `gtk_css_node_insert_after: assertion 'previous_sibling == NULL || previous_sibling->parent == parent' failed`
* IconView
  * Find why downloading data panics when called from newIconView.
* Cairo
  * Random crashes when cairo is used in the drawing area:
```
    gallery: ../cairo/src/cairo.c:523: cairo_destroy: Assertion `CAIRO_REFERENCE_COUNT_HAS_REFERENCE (&cr->ref_count)' failed.
    SIGABRT: abort
    PC=0x7f2595fbcd22 m=11 sigcode=18446744073709551610
    
    github.com/diamondburned/gotk4/pkg/gtk/v4._Cfunc_cairo_destroy(0x7f24d00093e0)
    	_cgo_gotypes.go:2980 +0x45 fp=0xc0000405c0 sp=0xc000040598 pc=0x7aed25
    github.com/diamondburned/gotk4/pkg/gtk/v4._gotk4_gtk4_DrawingAreaDrawFunc.func1.1(0xc000040620)
    	/github.com/diamondburned/gotk4/pkg@v0.0.0-20210919215506-2625db339437/gtk/v4/gtkdrawingarea.go:50 +0x53     fp=0xc0000405f8 sp=0xc0000405c0 pc=0x80cb13
    github.com/diamondburned/gotk4/pkg/gtk/v4._gotk4_gtk4_DrawingAreaDrawFunc.func1(0x178)
    	/github.com/diamondburned/gotk4/pkg@v0.0.0-20210919215506-2625db339437/gtk/v4/gtkdrawingarea.go:50 +0x19     fp=0xc000040610 sp=0xc0000405f8 pc=0x80ca99
    runtime.call16(0x0, 0xa26a78, 0x0, 0x0, 0x0, 0x0, 0xc0000406c0)
```