# ratata
Ratata (not to be mistaken with [Rattata](https://www.pokemon.com/us/pokedex/rattata)) is a Brazilian slang for spliting the bill. This is a Splitwise-like application that aims to simplify cost splitting between multiple people.

## Splitwise already exists, why Ratata?
I was looking for a simple way to split costs between people, for both short-term purposes (barbecues, group trips) and long-term (household stuff). Splitwise is fine, but I was looking for something that fits my needs better:
* Unlimited daily transactions (for free)
* [TBD] Web-based (no need to install apps)
* [TBD] No registration (for those trips with people that don't use such apps)
* [TBD] Cost breakdown reports

## Quick Start

## Usage

## Contributing
### Clone the repo
```bash
git clone https://github.com/bhashimoto/ratata@latest
cd ratata
```

### Set up environment variables
**.env**
```
DATABASE_URL=""
PORT=""
```

### Run migrations
```
./scripts/migrate.sh up
```

### Build application
```bash
go build
```

### Run application
```bash
./ratata
```


