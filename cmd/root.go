/*
Copyright © 2021 Barbé Creuly Bastien

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "dimpsled",
		Short: "Set a dim light for a connected PS controller",
		Long: `dimpsled takes a PS controller device, or any device that behave similarly,
and choose a random dim, non aggressive light and color, so that it does not burn the eyes of peoples in front of you.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			if red, green, blue, err := findRGBFiles(globalDevice); err != nil {
				panic(err)
			} else {
				if col, err := generateRGBcolor(palette); err != nil {
					panic(err)
				} else {
					r, g, b := col.RGB255()
					fmt.Printf("Red (%s): %d\nGreen (%s): %d\nBlue (%s): %d\n", red, r, green, g, blue, b)
					if err := setPSLEDColors(red, green, blue, col); err != nil {
						panic(err)
					}
				}
			}
		},
	}
	globalDevice string = ""
	palette      string = ""
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dimpsled.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rand.Seed(time.Now().UTC().UnixNano())
	rootCmd.Flags().StringVarP(&globalDevice, "device", "d", "", "Any device ending with @global, for wich there is a corresponding @red, @green and @blue")
	rootCmd.Flags().StringVarP(&palette, "palette", "p", "pastelle", "Color palette to use (pastelle)")
}

func findRGBFiles(dev string) (string, string, string, error) {
	if _, err := os.Stat(dev); err != nil {
		return "", "", "", err
	}

	// We need to cut the @global or :global from the given filename
	subPartParser := regexp.MustCompile(`(^.*[@:])global$`)
	subPart := subPartParser.FindStringSubmatch(dev)[1]

	// Now we create red green and blue devicepath
	format := "%s%s"
	red, green, blue := fmt.Sprintf(format, subPart, "red"), fmt.Sprintf(format, subPart, "green"), fmt.Sprintf(format, subPart, "blue")

	// We still need to check their existence or any errors
	for _, path := range []string{red, green, blue} {
		if _, err := os.Stat(path); err != nil {
			return "", "", "", err
		}
	}

	return red, green, blue, nil
}

func generateRGBcolor(pal string) (colorful.Color, error) {
	var paletteSettings colorful.SoftPaletteSettings
	var color colorful.Color
	switch pal {
	case "pastelle":
		paletteSettings = colorful.SoftPaletteSettings{
			CheckColor:  pastelle,
			Iterations:  70,
			ManySamples: true,
		}
	default:
		return color, fmt.Errorf("Unknown color palette %s", pal)
	}

	// Generate only one color from the palette generator
	if palette, err := colorful.SoftPaletteEx(1, paletteSettings); err != nil {
		return color, err
	} else {
		color = palette[0].Clamped()
	}

	return color, nil
}

func pastelle(l, a, b float64) bool {
	// This is the result of many tests, I wanted dark colors because the PS controller LED can be be very bright
	h, c, L := colorful.LabToHcl(l, a, b)
	return 100 < h && h < 300 && 0.1 < c && c < 0.2 && L < 0.1
}

func setPSLEDColors(red, green, blue string, color colorful.Color) error {
	r, g, b := color.RGB255()
	// Set red
	if err := writeColorToFile(red, r); err != nil {
		return err
	}
	if err := writeColorToFile(green, g); err != nil {
		return err
	}
	if err := writeColorToFile(blue, b); err != nil {
		return err
	}
	return nil
}

func writeColorToFile(filename string, color uint8) error {
	if _, err := os.Stat(filename); err != nil {
		return err
	} else {
		brightness := fmt.Sprintf("%s/brightness", filename)
		if file, err := os.OpenFile(brightness, os.O_WRONLY, os.ModeDevice|os.ModeCharDevice); err != nil {
			return err
		} else {
			if _, err := file.WriteString(fmt.Sprintf("%d", color)); err != nil {
				return err
			}
		}
	}
	return nil
}
