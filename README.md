# MarkovGenBot
MarkovGenBot is a bot for Telegram written in Go which implements a markovian generator with a single level of memory.



## Theory
A markovian text generator is a system for generating text based on one or more training texts.

The data structure consists of a dictionary where the keys are words and the values are word-integer dictionaries. Given a word w, its value contains words that followed w in the training, each with its number of occurrences.


### Training phase
Each training updates the data structure incrementing, for each pair of consecutive words, the number of occurrences of the latter word in the dictionary of the former.

Training can be easily repeated multiple times if using a mutable implementation: there is no need for a separate phase.

In this implementation, the first word of the training text is considered to be a follower of the empty word "", while the last word of the training text has "" as follower.


### Generating phase
Given a trained data structure, the generator builds a text word by word. Given a word w, the next generated word is chosen probabilistically among the ones from the dictionary of w, each candidate with its number of occurrences as weight.

In this implementation, the first word is considered to be "", and the generation stops when the next generated word is "".



## Implementation
The choice of using "" as first and last word, in conjunction of it being considered the first and last word of the trainings, implies that any generated text starts with a word that started a training text and ends with a word that ended one. This is not an issue in this context, as each message the bot receives is a separate training set, and actually makes the generated texts sound more realistic.

The data structure is basically implemented as a `map[string]map[string]int` for each chat. The system trains on each received text, and generates whenever either the `/generate` command is used or a user replies to a message sent my the bot. The training data is periodically written to persistent files and unloaded if not used for a predetermined time span, while it is obviously loaded if necessary and not already in memory.



## Usage
- Download and build:
```sh
go install -v github.com/sgorblex/MarkovGenBot@latest
```
- Insert an API token obtained from @BotFather in a file `api_key.txt`
- run `MarkovGenBot` (the binary can be found in `$GOPATH`):
Remember to make sure your bot is configured to have access to messages and groups.

You can restrict chats which can interact with the bot by putting a json list of allowed chat IDs in `whitelist.json`. For example:
```json
[420000,-69420]
```
If the file `whitelist.json` does not exist, the bot will operate in no-whitelist mode, i.e. will allow conversation with any chat.

You will find persistent data in `data/<chatID>.json`, for each chat ID with which the bot interacted.



## License
[MIT](LICENSE)
