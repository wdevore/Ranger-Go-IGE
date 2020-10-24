package color

const (
	// Transparent colorless
	Transparent = uint64(0x00000000)

	// White color
	White = uint64(0xffffffff)
	// Black color
	Black = uint64(0x000000ff)
	// Red color
	Red = uint64(0xff0000ff)
	// Green color
	Green = uint64(0x00ff00ff)
	// Blue color
	Blue = uint64(0x0000ffff)

	// DarkerGray color
	DarkerGray = uint64(0x404040ff)
	// DarkGray color
	DarkGray = uint64(0x505050ff)
	// Gray color
	Gray = uint64(0xaaaaaaff)
	// LightGray color
	LightGray = uint64(0x646464ff)
	// LighterGray color
	LighterGray = uint64(0x848484ff)
	// LightestGray color
	LightestGray = uint64(0xa4a4a4ff)
	// Silver color
	Silver = uint64(0xddddddff)

	// ----------------------------------------------------
	// Red hues
	// ----------------------------------------------------

	// Pink color
	Pink = uint64(0xFFC0CBff)
	// LightPink color
	LightPink = uint64(0xFFB6C1ff)
	// HotPink color
	HotPink = uint64(0xFF69B4ff)
	// DeepPink color
	DeepPink = uint64(0xFF1493ff)

	// ----------------------------------------------------
	// Yellow hues
	// ----------------------------------------------------

	// DarkOrange color
	DarkOrange = uint64(0xff7f00ff)
	// Orange color
	Orange = uint64(0xff8000ff)
	// LightOrange color
	LightOrange = uint64(0xffa000ff)
	// SoftOrange color
	SoftOrange = uint64(0xFF851BFF)
	// Yellow color
	Yellow = uint64(0xFFDC00FF)
	// YellowGreen color
	YellowGreen = uint64(0x9ACD32FF)
	// GoldYellow color
	GoldYellow = uint64(0xFFC800FF)
	// Peach color
	Peach = uint64(0xF1C6A7FF)

	// ----------------------------------------------------
	// Green hues
	// ----------------------------------------------------

	// SoftGreen color
	SoftGreen = uint64(0x2ECC40FF)
	// GreenYellow color
	GreenYellow = uint64(0xADFF2FFF)
	// Olive color
	Olive = uint64(0x3D9970FF)
	// Teal color
	Teal = uint64(0x39CCCCFF)
	// Lime color
	Lime = uint64(0x01FF70FF)

	// ----------------------------------------------------
	// Blue hues
	// ----------------------------------------------------

	// SoftBlue color
	SoftBlue = uint64(0x0074D9FF)
	// DarkBlue color
	DarkBlue = uint64(0x6D9DEBFF)
	// GreyBlue color
	GreyBlue = uint64(0x4864B4FF)
	// Navy color
	Navy = uint64(0x001f3fFF)
	// Aqua color
	Aqua = uint64(0x7FDBFFFF)
	// LightPurple color
	LightPurple = uint64(0xaaaaffFF)

	// ----------------------------------------------------
	// Brown hues
	// ----------------------------------------------------

	// Chocolate color
	Chocolate = uint64(0xD2691EFF)
	// Saddlebrown color
	Saddlebrown = uint64(0x8B4513FF)
	// Sienna color
	Sienna = uint64(0xA0522DFF)
	// Brown color
	Brown = uint64(0xA52A2AFF)
	// Brick color
	Brick = uint64(0xC54A4AFF)

	// ----------------------------------------------------
	// Pantones hues
	// https://www.pantone.com/color-finder#/pick?pantoneBook=pantoneSolidCoatedV3M2
	// ----------------------------------------------------

	// LightNavyBlue color
	LightNavyBlue = uint64(0x85B3D1FF)
	// PanSkin color
	PanSkin = uint64(0xfcc89bff)
	// PanPurple color
	PanPurple = uint64(0x8031a7ff)

	// ----------------------------------------------------
	// Pastels
	// https://www.schemecolor.com/pastel-color-tones.php
	// ----------------------------------------------------

	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^
	// Unvaryingly Simple
	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^

	// UnSiCreamBlue --
	UnSiCreamBlue = uint64(0xA7BAE7FF)
	// UnSiSkyBlue --
	UnSiSkyBlue = uint64(0xD0E5FAFF)
	// UnSiLightPink --
	UnSiLightPink = uint64(0xD9BBE6FF)
	// UnSiDarkPurple --
	UnSiDarkPurple = uint64(0x57466DFF)

	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^
	// Greens and Blues
	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^

	// GBsCackyGreen --
	GBsCackyGreen = uint64(0xAFD5AAFF)
	// GBsSunRiseBlue --
	GBsSunRiseBlue = uint64(0x83C6DDFF)
	// GBsCornerBlue --
	GBsCornerBlue = uint64(0x5DB1D1FF)

	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^
	// What is true?
	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^

	// WiTPink --
	WiTPink = uint64(0xFDB4C1FF)
	// WiTSkin --
	WiTSkin = uint64(0xFFDEDAFF)
	// WiTSkyLavender --
	WiTSkyLavender = uint64(0xC5C2DFFF)
	// WiTCornerPurple --
	WiTCornerPurple = uint64(0x9B94BEFF)
	// WiTPeachCake --
	WiTPeachCake = uint64(0xF9C8A0FF)
	// WiTRibbonRed --
	WiTRibbonRed = uint64(0xF9C8A0FF)

	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^
	// Cacky browns
	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^

	// CBTeaLeaf --
	CBTeaLeaf = uint64(0xB1C5C3FF)
	// CBPeach --
	CBPeach = uint64(0xE9E2D7FF)
	// CBClay --
	CBClay = uint64(0xE1B894FF)
	// CBBrown --
	CBBrown = uint64(0x875C36FF)
	// CBDirtBrown --
	CBDirtBrown = uint64(0x704523FF)

	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^
	// Marriage Approved
	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^

	// MASkyLime --
	MASkyLime = uint64(0xD1D693FF)
	// MACackyGreen --
	MACackyGreen = uint64(0xBABF5EFF)
	// MABurgandy --
	MABurgandy = uint64(0x602F44FF)

	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^
	// Brillants
	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^

	// BrSmolderRed --
	BrSmolderRed = uint64(0xC33C23FF)
	// BrCinderRed --
	BrCinderRed = uint64(0xDC453DFF)
	// BrCreamRed --
	BrCreamRed = uint64(0xFF6961FF)
	// BrPuffYellow --
	BrPuffYellow = uint64(0xFFFD96FF)
)
