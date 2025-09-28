
# Build Your Own URL Shortener

<https://codingchallenges.fyi/challenges/challenge-url-shortener/>

## Step 0

Brain storming and ideas for implementing the url shortener.
I think I should approach the app in this order.

1. Create the functions that will convert the urls.
2. Set up the Database to store and retrieve information.
3. Need to handle conflicts for stored information.
4. Create end points for functions in parts 1 & 2.
5. Serve a basic webpage that does parts 1 & 2.

- REST API
  - <https://www.jetbrains.com/guide/go/tutorials/rest_api_series/stdlib/>
  - Return suitable HTTP Status Code
- Where should the shortend URLs & full URLs be saved. (SQL?)
  - I think I'll use SQLlite. Otherwise Postgres.
  - <https://practicalgobook.net/posts/go-sqlite-no-cgo/#prerequisites>
  - <https://www.twilio.com/en-us/blog/developers/community/use-sqlite-go>
  - <https://medium.com/@peymaan.abedinpour/golang-crud-app-tutorial-step-by-step-guide-using-sqlite-a3ce08a4fc81>
- Hash function to create a short key
- HTML for the display
  - <https://go.dev/doc/articles/wiki/>
  - <https://medium.com/@kodeforce/build-a-simple-server-in-golang-that-serves-basic-html-6af32f8525ed>
  - <https://coderonfleek.medium.com/htmx-go-build-a-crud-app-with-golang-and-htmx-081383026466>

<https://www.programming-books.io/essential/go/>
<https://medium.com/homeday/collision-handling-in-our-url-shortener-service-6612e3e82eae>
<https://dev.to/mokiat/proper-http-shutdown-in-go-3fji>

### Step 0.1

Hash Function Chosen: FNV (Fowler-Noll-Vo) Hash

- If there's a conflict, I can use a larger hash (FNV-32, FNV-64, FNV-128).
- Maybe include time to make sure hash is unique?
- I should check if longURL exists & then shortURL.
- I could also pepper the text to avoid collisions
Had to also use url.QueryEscape just so encoded hashes are url compatible
This might be a good choice for this project. It's simple and fast.

### Step 0.2

Used SQLLite3 DB.

- Set up temp db for now
- Once ready I will commit
