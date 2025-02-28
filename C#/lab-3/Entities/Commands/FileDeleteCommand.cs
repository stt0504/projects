using System;
using System.Collections.ObjectModel;
using Itmo.ObjectOrientedProgramming.Lab4.Models;

namespace Itmo.ObjectOrientedProgramming.Lab4.Entities.Commands;

public class FileDeleteCommand : CommandBase
{
    public FileDeleteCommand(Collection<Flag> flags, string path)
        : base(flags)
    {
        Path = path;
    }

    public string Path { get; private set; }

    public override void Execute(FileSystem fileSystem)
    {
        if (fileSystem is null) throw new ArgumentNullException(nameof(fileSystem));
        if (fileSystem.CurrentDirectory is null) throw new ArgumentException("current directory is null");

        if (!Path.StartsWith(fileSystem.CurrentDirectory, StringComparison.CurrentCulture))
        {
            Path = System.IO.Path.Combine(fileSystem.CurrentDirectory, Path);
        }

        fileSystem.DeleteFile(Path);
    }
}