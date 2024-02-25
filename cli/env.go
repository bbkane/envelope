package cli

import (
	"errors"
	"fmt"
	"time"

	"github.com/alessio/shellescape"
	"go.bbkane.com/envelope/domain"
	"go.bbkane.com/envelope/tableprint"

	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/value/scalar"
)

func EnvCreateCmd() command.Command {
	return command.New(
		"Create an environment",
		envCreateRun,
		command.ExistingFlags(commonCreateFlagMap()),
		command.ExistingFlag("--name", envNameFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
	)
}

func envCreateRun(cmdCtx command.Context) error {

	// common create Flags
	comment := cmdCtx.Flags["--comment"].(string)
	createTime := cmdCtx.Flags["--create-time"].(time.Time)
	updateTime := cmdCtx.Flags["--update-time"].(time.Time)

	name := cmdCtx.Flags["--name"].(string)

	iesr, err := initEnvService(cmdCtx.Flags)
	if err != nil {
		return err
	}
	defer iesr.Cancel()

	env, err := iesr.EnvService.EnvCreate(iesr.Ctx, domain.EnvCreateArgs{
		Name:       name,
		Comment:    comment,
		CreateTime: createTime,
		UpdateTime: updateTime,
	})

	if err != nil {
		return fmt.Errorf("could not create env: %w", err)
	}

	fmt.Fprintf(cmdCtx.Stdout, "Created env: %s\n", env.Name)

	return nil
}

func EnvDeleteCmd() command.Command {
	return command.New(
		"Delete an environment and associated vars",
		envDeleteRun,
		command.ExistingFlag("--name", envNameFlag()),
		command.ExistingFlags(confirmFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
	)
}

func envDeleteRun(cmdCtx command.Context) error {

	name := cmdCtx.Flags["--name"].(string)
	confirm := cmdCtx.Flags["--confirm"].(bool)

	if confirm {
		keepGoing, err := askConfirm()
		if err != nil {
			panic(err)
		}
		if !keepGoing {
			return errors.New("unconfirmed change")
		}
	}

	iesr, err := initEnvService(cmdCtx.Flags)
	if err != nil {
		return err
	}
	defer iesr.Cancel()

	err = iesr.EnvService.EnvDelete(iesr.Ctx, name)

	if err != nil {
		return fmt.Errorf("could not delete env: %s: %w", name, err)
	}

	fmt.Fprintf(cmdCtx.Stdout, "deleted: %s\n", name)
	return nil
}

func EnvListCmd() command.Command {
	return command.New(
		"List environments",
		envListRun,
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.ExistingFlags(timeZoneFlagMap()),
	)
}

func envListRun(cmdCtx command.Context) error {

	timezone := cmdCtx.Flags["--timezone"].(string)

	iesr, err := initEnvService(cmdCtx.Flags)
	if err != nil {
		return err
	}
	defer iesr.Cancel()

	envs, err := iesr.EnvService.EnvList(iesr.Ctx)

	if err != nil {
		return err
	}
	tableprint.EnvList(cmdCtx.Stdout, envs, tableprint.Timezone(timezone))
	return nil
}

func EnvPrintScriptCmd() command.Command {
	return command.New(
		"Print export script",
		envPrintScriptRun,
		command.ExistingFlag("--name", envNameFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.Flag(
			"--no-env-no-problem",
			"Exit without an error if the environment doesn't exit. Useful when runnng envelop on chpwd",
			scalar.Bool(
				scalar.Default(false),
			),
			flag.Required(),
		),
		command.Flag(
			"--type",
			"Type of script",
			scalar.String(
				scalar.Choices("export", "unexport"),
				scalar.Default("export"),
			),
			flag.Required(),
		),
	)
}

func envPrintScriptRun(cmdCtx command.Context) error {
	name := cmdCtx.Flags["--name"].(string)
	scriptType := cmdCtx.Flags["--type"].(string)
	noEnvNoProblem := cmdCtx.Flags["--no-env-no-problem"].(bool)

	iesr, err := initEnvService(cmdCtx.Flags)
	if err != nil {
		return err
	}
	defer iesr.Cancel()

	envVars, err := iesr.EnvService.EnvLocalVarList(iesr.Ctx, name)
	if err != nil {
		if errors.Is(err, domain.ErrEnvNotFound) && noEnvNoProblem {
			return nil
		}
		return fmt.Errorf("could not list env vars: %s: %w", name, err)
	}

	switch scriptType {
	case "export":
		for _, ev := range envVars {
			fmt.Fprintf(cmdCtx.Stdout, "echo 'Adding:' %s;\n", shellescape.Quote(ev.Name))
			fmt.Fprintf(cmdCtx.Stdout, "export %s=%s;\n", shellescape.Quote(ev.Name), shellescape.Quote(ev.Value))
		}
	case "unexport":
		for _, ev := range envVars {
			fmt.Fprintf(cmdCtx.Stdout, "echo 'Removing:' %s;\n", shellescape.Quote(ev.Name))
			fmt.Fprintf(cmdCtx.Stdout, "unset %s;\n", shellescape.Quote(ev.Name))
		}
	default:
		return errors.New("Unimplemented --script-type: " + scriptType)

	}

	return nil
}

func EnvShowCmd() command.Command {
	return command.New(
		"Print environment details",
		envShowRun,
		command.ExistingFlag("--name", envNameFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.ExistingFlags(timeZoneFlagMap()),
	)
}

func envShowRun(cmdCtx command.Context) error {

	name := cmdCtx.Flags["--name"].(string)
	timezone := cmdCtx.Flags["--timezone"].(string)

	iesr, err := initEnvService(cmdCtx.Flags)
	if err != nil {
		return err
	}
	defer iesr.Cancel()

	env, err := iesr.EnvService.EnvShow(iesr.Ctx, name)
	if err != nil {
		return fmt.Errorf("could not show env: %s: %w", name, err)
	}

	localvars, err := iesr.EnvService.EnvLocalVarList(iesr.Ctx, name)
	if err != nil {
		return err
	}

	refs, referencedVars, err := iesr.EnvService.EnvRefList(iesr.Ctx, name)
	if err != nil {
		return err
	}

	tableprint.EnvShowRun(cmdCtx.Stdout, *env, localvars, refs, referencedVars, tableprint.Timezone(timezone))
	return nil
}

func EnvUpdateCmd() command.Command {
	return command.New(
		"Update an environment",
		envUpdateRun,
		command.ExistingFlags(commonUpdateFlags()),
		command.ExistingFlag("--name", envNameFlag()),
		command.ExistingFlags(timeoutFlagMap()),
		command.ExistingFlags(sqliteDSNFlagMap()),
		command.ExistingFlags(confirmFlag()),
	)
}

func envUpdateRun(cmdCtx command.Context) error {

	// common update flags
	comment := ptrFromMap[string](cmdCtx.Flags, "--comment")
	createTime := ptrFromMap[time.Time](cmdCtx.Flags, "--create-time")
	newName := ptrFromMap[string](cmdCtx.Flags, "--new-name")
	updateTime := ptrFromMap[time.Time](cmdCtx.Flags, "--update-time")

	confirm := cmdCtx.Flags["--confirm"].(bool)
	name := cmdCtx.Flags["--name"].(string)

	if confirm {
		keepGoing, err := askConfirm()
		if err != nil {
			panic(err)
		}
		if !keepGoing {
			return errors.New("unconfirmed change")
		}
	}

	iesr, err := initEnvService(cmdCtx.Flags)
	if err != nil {
		return err
	}
	defer iesr.Cancel()

	err = iesr.EnvService.EnvUpdate(iesr.Ctx, name, domain.EnvUpdateArgs{
		Comment:    comment,
		CreateTime: createTime,
		NewName:    newName,
		UpdateTime: updateTime,
	})

	if err != nil {
		return fmt.Errorf("could not update env: %w", err)
	}
	finalName := name
	if newName != nil {
		finalName = *newName
	}
	fmt.Fprintln(cmdCtx.Stdout, "updated env:", finalName)
	return nil
}
