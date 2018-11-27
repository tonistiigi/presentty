package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	dockercommand "github.com/tonistiigi/presentty/dockercommand"
	"github.com/yudai/gotty/server"
	"github.com/yudai/gotty/utils"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var configFile string
	flag.StringVar(&configFile, "config", "config.toml", "Config file")
	flag.Parse()

	cfg, _, err := LoadFile(configFile)
	if err != nil {
		return err
	}

	if len(cfg.Demos) == 0 {
		return errors.New("no demos in config")
	}

	m, err := provisionDemos(filepath.Dir(configFile), cfg.Demos)
	if err != nil {
		return err
	}

	log.Printf("provisioned demos: %+v", m)
	defer func() {
		cleanup(m)
	}()

	factory, err := dockercommand.NewFactory(m, &dockercommand.Options{})
	if err != nil {
		return err
	}

	appOptions := &server.Options{}

	if err := utils.ApplyDefaultValues(appOptions); err != nil {
		return err
	}
	if cfg.Light {
		colorPaletteOverides := []string{
			"#073642", /*  0: black    */
			"#dc322f", /*  1: red      */
			"#859900", /*  2: green    */
			"#b58900", /*  3: yellow   */
			"#268bd2", /*  4: blue     */
			"#d33682", /*  5: magenta  */
			"#2aa198", /*  6: cyan     */
			"#eee8d5", /*  7: white    */
			"#002b36", /*  8: brblack  */
			"#cb4b16", /*  9: brred    */
			"#586e75", /* 10: brgreen  */
			"#506067", /* 11: bryellow */
			"#839496", /* 12: brblue   */
			"#6c71c4", /* 13: brmagenta*/
			"#93a1a1", /* 14: brcyan   */
			"#fdf6e3", /* 15: brwhite  */
		}

		var cpo []*string
		for _, v := range colorPaletteOverides {
			vv := v
			cpo = append(cpo, &vv)
		}

		appOptions.Preferences = &server.HtermPrefernces{
			CtrlPlusMinusZeroZoom: true,
			BackgroundColor:       colorPaletteOverides[15],
			ForegroundColor:       colorPaletteOverides[11],
			CursorColor:           "rgba(101, 123, 131, 0.4)",
			// FontSize:              18,
			ColorPaletteOverrides: cpo,
			// EnableBoldAsBright:    false,
		}
	} else {
		appOptions.Preferences = &server.HtermPrefernces{
			CtrlPlusMinusZeroZoom: true,
			// EnableBoldAsBright:    false,
		}
	}
	if cfg.Size != 0 {
		appOptions.Preferences.FontSize = cfg.Size
	}
	appOptions.PermitWrite = true

	// utils.ApplyFlags(cliFlags, flagMappings, c, appOptions, nil)

	srv, err := server.New(factory, appOptions)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	gCtx, gCancel := context.WithCancel(context.Background())

	errs := make(chan error, 1)
	go func() {
		errs <- srv.Run(ctx, server.WithGracefullContext(gCtx))
	}()
	err = waitSignals(errs, cancel, gCancel)

	if err != nil && err != context.Canceled {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	return nil
}

func waitSignals(errs chan error, cancel context.CancelFunc, gracefullCancel context.CancelFunc) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	select {
	case err := <-errs:
		return err

	case s := <-sigChan:
		switch s {
		case syscall.SIGINT:
			gracefullCancel()
			fmt.Println("C-C to force close")
			select {
			case err := <-errs:
				return err
			case <-sigChan:
				fmt.Println("Force closing...")
				cancel()
				return <-errs
			}
		default:
			cancel()
			return <-errs
		}
	}
}
