# 27. configure opens a file every other command refuses

- Status: accepted
- Date: 2026-09-13
- Part of: socrates-review configure

## Decision

A configuration file socrates-review cannot read is a refusal.
That is right for every command that acts on it: a review running on a configuration nobody can read is exactly the danger the refusal exists for, and running on defaults while a file says otherwise would be worse than stopping.

`socrates-review configure` is the one exception.
Its whole purpose is to spare a person the editor, and the moment they most need that is the moment they have broken the file by hand.
Refusing there sends them to the editor with a `yaml:` error for a map, which is the failure the command was written to remove.

So `configure` opens the file anyway.
It reads the file as nested regions of lines, one per key, and offers each region to the loader on its own alongside the key lines above it.
It names the deepest region the loader will not take, so a `temperature` nobody has under `roles.reviewer` is reported as `roles.reviewer.temperature` rather than as `roles`: the line it starts on, the text of that line, and the loader's own account of what is wrong with it.
The only thing it offers is resetting that setting to its default, which is removing it.
It guesses at nothing, it rewrites nothing else, and every other key and comment in the file is left where it was.
Once the offending block is gone the file parses and the menu carries on.

Reading the file as regions rather than trusting the parser's own line number is what makes the reported line the one that is actually wrong.
An unclosed flow sequence is reported by YAML at the end of the file, and by `configure` at the line that opened it.

Removing the named region goes through the same write every other change goes through where the file parses as a document at all, so an emptied block goes with its last setting and every comment stays where it was.
Where no single region accounts for the failure, `configure` refuses like every other command rather than offering to remove something it cannot name.

## Rejected

- refusing like every other command, which is consistent and leaves the person with the editor and a parser error.
- reading the parser's error to work out which key is wrong, which is matching on a dependency's message text and reports the wrong line for a whole class of mistakes.
- reading only the top-level blocks, which is a third of the code and names `roles` where the thing that is wrong is one line three levels down, taking two settings that were fine with it.
- dropping what cannot be parsed and carrying on, which silently throws away something a person wrote.
- rewriting the file from the settings that did parse, which is the same silent loss wearing a helpful face.
