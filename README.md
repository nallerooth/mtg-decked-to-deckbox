# MTG Decked To Deckbox
A simple tool for converting CSV files exported from 
[Decked Builder](http://www.deckedbuilder.com/) to a CSV format that can be 
imported into [Deckbox](https://deckbox.org).

## Building
```
cd mtg-decked-to-deckbox
go install
```

## Usage
Given that you've added `$GOPATH/bin` to your $PATH, you should be able to
run the following command:
```
mtg-decked-to-deckbox <CSV file>
```

The result will be a new CSV file, named `converted.csv`, located in the 
current working directory.

### Notice
This is a hack. While it seems to work, it's not ready for production use
(whatever that may mean). 

### Known bugs and limitations
- It seems that the Decked Builder app exports the back of some 2-sided cards, 
which breaks the import at Deckbox. Those cards have been added to a static 
blacklist. This should really be moved into a more dynamic solution.
 
