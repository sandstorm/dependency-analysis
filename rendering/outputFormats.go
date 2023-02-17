package rendering

type OutputFormat struct {
	Label       string
	FileEnding  string
	FlagValue   string
	Description string
}

// Provides the OutputFormat for the given flag if present, nil otherwise
func GetOutputFormatByFlagValue(flagValue string) *OutputFormat {
	for _, format := range SupportedOutputFormats {
		if format.FlagValue == flagValue {
			return &format
		}
	}
	return nil
}

// copied from https://graphviz.org/docs/outputs/
var SupportedOutputFormats = []OutputFormat{
	{
		Label:       "BMP",
		FileEnding:  "bmp",
		FlagValue:   "bmp",
		Description: "Windows Bitmap",
	},
	{
		Label:       "CGImage",
		FileEnding:  "cgimage",
		FlagValue:   "cgimage	",
		Description: "Apple Core Graphics",
	},
	{
		Label:       "DOT",
		FileEnding:  "canon",
		FlagValue:   "canon",
		Description: "Graphviz Language",
	},
	{
		Label:       "DOT",
		FileEnding:  "dot",
		FlagValue:   "dot",
		Description: "Graphviz Language",
	},
	{
		Label:       "DOT",
		FileEnding:  "gv",
		FlagValue:   "gv",
		Description: "Graphviz Language",
	},
	{
		Label:       "DOT",
		FileEnding:  "xdot",
		FlagValue:   "xdot",
		Description: "Graphviz Language",
	},
	{
		Label:       "DOT",
		FileEnding:  "xdot",
		FlagValue:   "xdot1.2",
		Description: "Graphviz Language",
	},
	{
		Label:       "DOT",
		FileEnding:  "xdot",
		FlagValue:   "xdot1.4",
		Description: "Graphviz Language",
	},
	{
		Label:       "EPS",
		FileEnding:  "eps",
		FlagValue:   "eps",
		Description: "Encapsulated PostScript",
	},
	{
		Label:       "EXR",
		FileEnding:  "exr",
		FlagValue:   "exr",
		Description: "OpenEXR",
	},
	{
		Label:       "FIG",
		FileEnding:  "fig",
		FlagValue:   "fig",
		Description: "Xfig",
	},
	{
		Label:       "GD",
		FileEnding:  "gd",
		FlagValue:   "gd",
		Description: "LibGD",
	},
	{
		Label:       "GD2",
		FileEnding:  "gd2",
		FlagValue:   "gd2",
		Description: "LibGD",
	},
	{
		Label:       "GIF",
		FileEnding:  "gif",
		FlagValue:   "gif",
		Description: "Graphics Interchange Format",
	},
	{
		Label:       "GTK",
		FileEnding:  "gtk",
		FlagValue:   "gtk",
		Description: "Formerly GTK+ / GIMP ToolKit",
	},
	{
		Label:       "ICO",
		FileEnding:  "ico",
		FlagValue:   "ico",
		Description: "Windows Icon",
	},
	{
		Label:       "Image Map: Client-side",
		FileEnding:  "cmap",
		FlagValue:   "cmap",
		Description: "",
	},
	{
		Label:       "Image Map: Server-side",
		FileEnding:  "ismap",
		FlagValue:   "ismap",
		Description: "",
	},
	{
		Label:       "Image Map: Server-side and client-side",
		FileEnding:  "imap",
		FlagValue:   "imap",
		Description: "",
	},
	{
		Label:       "Image Map: Server-side and client-side",
		FileEnding:  "cmapx",
		FlagValue:   "cmapx",
		Description: "",
	},
	{
		Label:       "Image Map: Server-side and client-side",
		FileEnding:  "imap_np",
		FlagValue:   "imap_np",
		Description: "",
	},
	{
		Label:       "Image Map: Server-side and client-side",
		FileEnding:  "cmapx_np",
		FlagValue:   "cmapx_np",
		Description: "",
	},
	{
		Label:       "JPEG",
		FileEnding:  "jpg",
		FlagValue:   "jpg",
		Description: "Joint Photographic Experts Group",
	},
	{
		Label:       "JPEG",
		FileEnding:  "jpeg",
		FlagValue:   "jpeg",
		Description: "Joint Photographic Experts Group",
	},
	{
		Label:       "JPEG",
		FileEnding:  "jpe",
		FlagValue:   "jpe",
		Description: "Joint Photographic Experts Group",
	},
	{
		Label:       "JPEG 2000",
		FileEnding:  "jp2",
		FlagValue:   "jp2",
		Description: "",
	},
	{
		Label:       "JSON",
		FileEnding:  "json",
		FlagValue:   "json",
		Description: "JavaScript Object Notation",
	},
	{
		Label:       "JSON",
		FileEnding:  "json0",
		FlagValue:   "json0",
		Description: "JavaScript Object Notation",
	},
	{
		Label:       "JSON",
		FileEnding:  "dot_json",
		FlagValue:   "dot_json",
		Description: "JavaScript Object Notation",
	},
	{
		Label:       "JSON",
		FileEnding:  "xdot_json",
		FlagValue:   "xdot_json",
		Description: "JavaScript Object Notation",
	},
	{
		Label:       "PDF",
		FileEnding:  "pdf",
		FlagValue:   "pdf",
		Description: "Portable Document Format",
	},
	{
		Label:       "PIC",
		FileEnding:  "pic",
		FlagValue:   "pic",
		Description: "Brian Kernighan's Diagram Language",
	},
	{
		Label:       "PICT",
		FileEnding:  "pct",
		FlagValue:   "pct",
		Description: "Apple PICT",
	},
	{
		Label:       "PICT",
		FileEnding:  "pict",
		FlagValue:   "pict",
		Description: "Apple PICT",
	},
	{
		Label:       "Plain Text",
		FileEnding:  "plain",
		FlagValue:   "plain",
		Description: "Simple, line-based language",
	},
	{
		Label:       "Plain Text",
		FileEnding:  "plain-ext",
		FlagValue:   "plain-ext",
		Description: "Simple, line-based language",
	},
	{
		Label:       "PNG",
		FileEnding:  "png",
		FlagValue:   "png",
		Description: "Portable Network Graphics",
	},
	{
		Label:       "POV-Ray",
		FileEnding:  "pov",
		FlagValue:   "pov",
		Description: "Persistence of Vision Raytracer (prototype)",
	},
	{
		Label:       "PS",
		FileEnding:  "ps",
		FlagValue:   "ps",
		Description: "Adobe PostScript",
	},
	{
		Label:       "PS/PDF",
		FileEnding:  "ps2",
		FlagValue:   "ps2",
		Description: "Adobe PostScript for Portable Document Format",
	},
	{
		Label:       "PSD",
		FileEnding:  "psd",
		FlagValue:   "psd",
		Description: "Photoshop",
	},
	{
		Label:       "SGI",
		FileEnding:  "sgi",
		FlagValue:   "sgi",
		Description: "Silicon Graphics Image",
	},
	{
		Label:       "SVG",
		FileEnding:  "svg",
		FlagValue:   "svg",
		Description: "Scalable Vector Graphics",
	},
	{
		Label:       "SVG",
		FileEnding:  "svgz",
		FlagValue:   "svgz",
		Description: "Scalable Vector Graphics",
	},
	{
		Label:       "TGA",
		FileEnding:  "tga",
		FlagValue:   "tga",
		Description: "Truevision TARGA",
	},
	{
		Label:       "TIFF",
		FileEnding:  "tif",
		FlagValue:   "tif",
		Description: "Tag Image File Format",
	},
	{
		Label:       "TIFF",
		FileEnding:  "tiff",
		FlagValue:   "tiff",
		Description: "Tag Image File Format",
	},
	{
		Label:       "Tk",
		FileEnding:  "tk",
		FlagValue:   "tk",
		Description: "Tcl/Tk",
	},
	{
		Label:       "VML",
		FileEnding:  "vm",
		FlagValue:   "vm",
		Description: "Vector Markup Language",
	},
	{
		Label:       "VML",
		FileEnding:  "vmlz",
		FlagValue:   "vmlz",
		Description: "Vector Markup Language",
	},
	{
		Label:       "VRML",
		FileEnding:  "vrml",
		FlagValue:   "vrml",
		Description: "Virtual Reality Modeling Language",
	},
	{
		Label:       "WBMP",
		FileEnding:  "wbmp",
		FlagValue:   "wbmp",
		Description: "Wireless Bitmap",
	},
	{
		Label:       "WebP",
		FileEnding:  "webp",
		FlagValue:   "webp",
		Description: "WebP",
	},
	{
		Label:       "X11",
		FileEnding:  "xlib",
		FlagValue:   "xlib",
		Description: "X11 Window",
	},
	{
		Label:       "X11",
		FileEnding:  "x11",
		FlagValue:   "x11",
		Description: "X11 Window",
	},
}
