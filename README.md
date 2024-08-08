# Ratata
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![HTML](https://img.shields.io/badge/HTML-%23E34F26.svg?logo=html5&logoColor=white) ![HTMX](https://img.shields.io/badge/HTMX-36C?logo=htmx&logoColor=fff) ![SQLite](https://img.shields.io/badge/SQLite-%2307405e.svg?logo=sqlite&logoColor=white)]
Ratata (not to be mistaken with [Rattata](https://www.pokemon.com/us/pokedex/rattata)) is a Brazilian slang for spliting the bill. This is a Splitwise-like application that aims to simplify cost splitting between multiple people.

## Splitwise already exists, why Ratata?
I was looking for a simple way to split costs between people, for both short-term purposes (barbecues, group trips) and long-term (household stuff). Splitwise is fine, but I was looking for something that fits my needs better:
* Unlimited daily transactions
* Web-based (no need to install apps)
* [TBD] No registration (for those trips with people that don't use such apps)
* [TBD] Cost breakdown reports


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


