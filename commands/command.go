package commands

import "github.com/bwmarrin/discordgo"

type CommandContext struct {
	Session     *discordgo.Session
	Message     *discordgo.MessageCreate
	Interaction *discordgo.InteractionCreate
	Args        []string
}

type Command struct {
	Name        string
	Aliases     []string
	Description string
	Options     []*discordgo.ApplicationCommandOption

	Execute            func(context *CommandContext) error
	ExecuteInteraction func(context *CommandContext) error
}

type Registry struct {
	commands map[string]*Command
}

func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]*Command),
	}
}

func (registry *Registry) Register(command *Command) {
	registry.commands[command.Name] = command

	for _, alias := range command.Aliases {
		registry.commands[alias] = command
	}
}

func (registry *Registry) Find(name string) *Command {
	return registry.commands[name]
}

func (registry *Registry) ApplicationCommands() []*discordgo.ApplicationCommand {
	commands := make([]*discordgo.ApplicationCommand, 0)
	registered := make(map[*Command]struct{})

	for _, command := range registry.commands {
		if command.ExecuteInteraction == nil {
			continue
		}

		if _, exists := registered[command]; exists {
			continue
		}

		registered[command] = struct{}{}

		commands = append(
			commands,
			&discordgo.ApplicationCommand{
				Name:        command.Name,
				Description: command.Description,
				Options:     command.Options,
			},
		)
	}

	return commands
}
