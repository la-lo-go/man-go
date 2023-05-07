# Man-go

## Work in progress, currently arround 50% until 1.0 release
|                | **GetMangas** | **GetMangaPage** | **GetChapter** |
| :------------: | :-----------: | :--------------: | :------------: |
|  **In Manga**  |       ✅       |        ✅         |       ❌        |
|    **Nyaa**    |       ✅       |        ✅         |       ❌        |
|  **Manga MX**  |       ✅       |        ❌         |       ❌        |
| **TuMangaNet** |       ✅       |        ❌         |       ❌        |

![](https://miro.medium.com/max/1100/0*QamUdvAd_-mDNIlD.gif)

## Run the API
1. Clone the repository
    ``` bash
    git clone https://github.com/la-lo-go/man-go
    ```
2. Create a new .env file in the root of the repository with the following scheme:
    ``` env
    # API
    API_IP="0.0.0.0"
    API_PORT="7070"

    # DB
    DB_NAME="mango.sqlite3"
    ```
3. Compile and run
    ``` bash
    go run main.go 
    ```

## API responses
### Search
#### Query
```
/busqueda?search=[search_string]&max=[int (optional)]
```
#### Response
```json
Example: /busqueda?search=jujutsu&max=1

[
    {
        "name": "jujutsu-kaisen",
        "site": "InManga",
        "link": "https://inmanga.com/ver/manga/Jujutsu-Kaisen/d88692a5-c341-47fc-8e39-da11a8fdee82",
        "chapters_number": 0,
        "cover": "https://pack-yak.intomanga.com/thumbnails/manga/Jujutsu-Kaisen/d88692a5-c341-47fc-8e39-da11a8fdee82"
    },
    {
        "name": "jujutsu-kaisen",
        "site": "Nyaa",
        "link": "https://manganyaa.com/jujutsu-kaisen/leer-online-gratis-espanol",
        "chapters_number": 99,
        "cover": "https://content.manganyaa.com/file/mnyaaa/jujutsu-kaisen/description/1.jpg"
    },
    {
        "name": "jujutsu-kaisen",
        "site": "TuManga.net",
        "link": "https://tumanga.net/manga/jujutsu-kaisen",
        "chapters_number": 99,
        "cover": ""
    }
]
```

### Manga Page
#### Query
```
/manga/[manga_name]/[site]
```
#### Response
```json
Example: /manga/jujutsu-kaisen/nyaa

{
    "name": "Jujutsu Kaisen",
    "site": "Nyaa",
    "cover": "https://content.manganyaa.com/file/mnyaaa/jujutsu-kaisen/description/1.jpg",
    "chaptersNumber": 186,
    "chaptersListed": [
        {
            "number": 186,
            "linkOriginal": "https://manganyaa.com/jujutsu-kaisen/leer-online-gratis-espanol/capitulo/186"
        },
        {
            "number": 185,
            "linkOriginal": "https://manganyaa.com/jujutsu-kaisen/leer-online-gratis-espanol/capitulo/185"
        },
        ...
        {
            "number": 2,
            "linkOriginal": "https://manganyaa.com/jujutsu-kaisen/leer-online-gratis-espanol/capitulo/2"
        },
        {
            "number": 1,
            "linkOriginal": "https://manganyaa.com/jujutsu-kaisen/leer-online-gratis-espanol/capitulo/1"
        }
    ]
}
```

## Cache
- `API_searchCache.json`: all the results of the searches with the date.
  - The search is exact: the results are returned.
  - The has a key that is included in the search: the results are filtered.
  - Else: do a full search
``` json
[
    {
        "search": "shingeki no",
        "date": "2022-08-08T21:00:10.7641581+02:00",
        "response": [Search response]
    },
    {
        "search": "isekai",
        "date": "2022-08-10T16:30:51.9044505+02:00",
        "response": [Search response]
    }
]
```

- `API_linksMangas.json`: Group all the results with the same name.
``` json
"jujutsu-kaisen": {      <- search string, map key
    "name": "jujutsu-kaisen",
    "sitesLinks": [
      {
        "site": "InManga",
        "link": "https://inmanga.com/ver/manga/Jujutsu-Kaisen/d88692a5-c341-47fc-8e39-da11a8fdee82"
      },
      {
        "site": "Nyaa",
        "link": "https://manganyaa.com/jujutsu-kaisen/leer-online-gratis-espanol"
      },
      {
        "site": "TuManga.net",
        "link": "https://tumanga.net/manga/jujutsu-kaisen"
      }
    ]
}
```
- `API_linksMangas.json`: JSON returned by the Nyaa API with all the mangas listed in the page. Has a 12 hour expire time.


## File structure
```
man-go
├───endpoints           <- endpoints functions
├───gorm                <- gorm database functions
├───general_functions   <- handy custom functions
├───json                <- cached information
├───models              <- api's structs
├───routers             <- gin-gonic router
└───sites               <- sites using the interface
    ├───inManga
    ├───mangaMx
    ├───nyaa
    └───tuMangaNet
```


## License
[Mozilla Public License 2.0](https://www.mozilla.org/en-US/MPL/2.0/)

[![License: MPL 2.0](https://img.shields.io/badge/License-MPL_2.0-brightgreen.svg)](https://www.mozilla.org/en-US/MPL/2.0/)

