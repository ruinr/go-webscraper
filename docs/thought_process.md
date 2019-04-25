Development Thought Process
----
## Requirement Review
### Acceptance Criteria
- [x] Take Amazon Product ASIN as request variable
- [x] Fetch category, rank & product dimensions of the product
- [x] Store the data in some sort of database
- [x] Display the data on the front-end

### Questions need to be answered:
#### 1. What were the biggest challenges you faced in writing the challenge?
  The biggest challenge is probably to think about how to get this challenge done properly based on requirement and information I know. It means a lot of researching, analyzing, learning, and planning for data structure, technologies, tools, languages, libraries/gems/packages, etc. It is always time-consuming and challenging to make the right decision on these.

  The biggest technical challenge is probably writing the scraping service to get product data. It is mainly because I don't often do web crawler or scraper. In order to make this script work for most cases, I had to go back and forth to
  - check the HTML source
  - change the script
  - test the script
  - check scaped data
  - clean scaped data
  - ...

  And I haven't figure out a proper way to test the scraping script yet. The only way I can think of right now is to download the HTML source file for the different layout I have seen, and test against them.
  
  Ok. So...after some setbacks by bot detection, now I think the biggest challenge is avoid bot detection. This is actually the key for all web scraping. Stable is first, then performance. Lesson learned in a hard way! I researched some possible solutions to lower the chance of being detected.
  1. Proxy Switcher
  2. Slow down
  3. Change user-agent
  4. Cookies
  5. Random request header
  
  I added random delay for slowing down, set cookies, and randomize user-agent. I haven't figured out how to switch proxy on heroku and how to correctly randomize request header yet.

#### 2. Can you explain your thought process on how you solved the problem, and what iterations you went through to get to this solution?
  Almost all of my thought processes are documented below. There are two problems I did not document.
  - HTML entities inside a scraped string
    1. I tried HTML.UnescapeString(), but it doesn't convert everything
    2. Then I googled, and it turns out HTML.UnescapeString() is not 100% correct
    3. So, I add string replacement for some special cases

  - The Live demo doesn't work right before I want to submit the challenge, but it was working the day before
    1. I didn't realize that I didn't set user agent for the scraper at first, so I was checking what changes after the working version.
    2. After checking diffs, I didn't find anything weird or wrong. Then did some tests locally find out it works sometimes.
    3. Then I checked the on error response from Amazon. I found it sometimes returns "something wrong on our end please try again" and sometimes returns bot detected message. Then I realized that I didn't set user-agent for scraper

#### 3. If you had to scale this application to handle millions of people adding in ASIN's in a given day, what considerations would you make?
  1. Load balancer for both frontend and backend, and more copies of the application
  2. Improve caching solution and caching size
  3. Find the slow parts and improve them, which need to do benchmark

#### 4. Why did you choose the technologies you used?
  All the choices are explained below including approach, language, framework, database and API structure.

## Plan & Analysis
### 1. Amazon's Product API

  There is nothing better than using Amazon's own API to retrieve its product information.
  The note in the tech challenge states that "The API isn't available and you will need to figure out an alternative method."
  I wasn't really sure what the note means, and I didn't believe Amazon doesn't have the API endpoint to retrieve basic information like these. I briefly went through Amazon's Product Advertising API documentation and found that all the product information required by the challenge exist.

  Then asked Veerprit to clarify this note for me. It turns out that getting access to the API could take one week.

  It seems like my website needs to be reviewed first, and there is no developer access for playing around.

  https://associates.amazon.ca/assoc_credentials/home

  ![Alt Text](https://beemapp2.s3.amazonaws.com/f6164105-f552-48a5-af84-8da31eb85d82.jpg)

### 2. API is not available, what's the next best?

  Web Scraping. Cannot really think of any other approaches that don't require third-party resources.

### 3. What is the language I know best for Web Scraping?

  Python: Probably the first choice for many and the most popular one. However, I am still learning Python.
  
  PHP: I haven't really developed in PHP for a while, so probably not a good choice to do this challenge.
  
  Ruby vs. Golang: I choose Golang

  - Pros:
    * Famous for simplicity and productivity
    * Performance is awesome because of garbage, concurrency, compiled to machine code.
    * Type Safe
    * Microservices structure, easily scalable
    * First-class support for Protocol Buffers and gRPC
    * Has a robust built-in testing library
    * Has one of the fastest growing community
    * Not much recommendation for using Ruby to do web scraping
    * I am more familiar with how to set up a whole project better in Go
  - Cons:
    * Error Handling
    * Less flexible compared to Ruby and other languages
    * Lack of some basic built-in functions that exist in Ruby and other languages

### 4. Does Go has an open-source web scraping framework?

  Colly: https://github.com/gocolly/colly

  GoQuery: https://github.com/PuerkitoBio/goquery

  Briefly reviewed two frameworks, and decided to use Colly. Compared to GoQuery:

  - Pros:
    * Well documented
    * More examples to follow
    * Faster
    * Automatic cookie and session handling
    * Community is more active
  - Cons:
    * DOM traversal
    * GoQuery allows more complex scraping

  Read the blog post about Colly vs. Goquery

  [Scraping the Web in Golang with Colly and Goquery](https://benjamincongdon.me/blog/2018/03/01/Scraping-the-Web-in-Golang-with-Colly-and-Goquery), by Benjamin Congdon

  After checking out Amazon's product pages, I think this project doesn't need much complex scraping. Therefore, Colly is a better choice for performance, and it requires less time to get familiar with.

### 5. Data Structure
  Based on a brief observation of random products on both Amazon.com and Amazon.ca, I found the following:
  - A product can have one main category and multiple levels of subcategories
  - A product can have more than one dimension description and separated by ";" (https://www.amazon.ca/dp/B004QWYCVG)
  - A product can have no dimension (https://www.amazon.com/dp/B07FSH5L52)
  - A product can have one rank in its main category and multiple ranks in different subcategories
  - A product can have no rank info (https://www.amazon.ca/dp/B004QWYCVG)
  - A product can have only one ASIN
  - The same product seems to share one ASIN on both Amazon.com and Amazon.ca
  - units of measurement that used to describe product dimensions are different depends on countries or domains.

### 6. Database Choice
  After reviewing the requirements, I found some details like how and where this information is going to be used are not included, except displaying data to frontend.

  Based on the briefly observed product data structure and project requirements, NoSQL Document Database is probably the best approach because of flexible schema and scalability. However, choosing a database for a project is never easy. Relational databases still have their advantages like handling transactional data, well supported BI & analytics tools, etc. Without knowing all the details, it is very hard to choose either.

  However, I'd like to take the initiative to use Redis for this project's first phase instead. Here are some reasons:
  - Required data is not very complex, and a key-value NoSQL database is simple and enough for this project.
  - For this particular project, data storage size is not a concern, at least not at the moment.
  - Redis is lightweight, scalable and very fast
  - Redis is an in-memory database, so it means no runtime dependency.
  - Redis can also be used as a caching solution

### 7. API Design
  gRPC vs. REST

  [REST is not the Best for Micro-Services GRPC and Docker makes a compelling case](https://hackernoon.com/rest-in-peace-grpc-for-micro-service-and-grpc-for-the-web-a-how-to-908cc05e1083), by Alex Punnen

  [REST vs. gRPC: Battle of the APIs](https://code.tutsplus.com/tutorials/rest-vs-grpc-battle-of-the-apis--cms-30711), by Gigi Sayfan

  - Gain
    - Performance on almost every aspect
    - Experience on gRPC and Protocol Buffer
    - Credit for trying something new, maybe?

  - Loss
    - Code quality due to unfamiliarity
    - Compared to REST, gRPC is relatively new and not widely used at the moment
    - More time needs to spend on learning

  For performance and learning, I decide to embrace challenges and be adventurous on this one to use gRPC.

  Based on what I know about gRPC, it uses TCP, it means I need to add an HTTP gateway for the frontend to communicate. And it exists in [grpc-ecosystem](https://github.com/grpc-ecosystem)

### 8. Layout
  - Frontend + Frontend Caching: ReactJs + Nginx
  - REST gateway: REST endpoint to allow frontend to communicate with gRPC
  - gRPC server
  - Redis as backend caching, probably TTL only
  - Redis as an in-memory database for storing data

LET'S ROCK!!!!

## Development Phases
### First Phase
  Make It Work!
  - ~~A working scaper service that takes ASIN and returns required product info in struct successfully for all the cases mentioned Section#5 Data Structure.~~
  - ~~A gRPC server with a RPC handler that calls scraper service.~~
  - ~~A REST gateway allows REST API calls to connect to gRPC server.~~
  - ~~Scaped product info can be stored into Redis as in-memory database~~
  - ~~Scaped product info can be returned as JSON to API call~~
  - ~~A basic cache for storing scraped product~~
  - ~~Basic testing for services~~
  - ~~A successfully built docker image~~
  - ~~A working live demo for API call~~

### Phase two
  Improve it!
  - ~~Add proper logging~~
  - ~~Improve scaped content~~
  - Add certificate-based authentication to gRPC server
  - Add token-based authentication to REST gateway
  - Add middleware authentication
  - Add more functionality like getting searched item list or category list, take domain as variable and scrap for other Amazon domains
  - Add more tests
  - Benchmark for backend
  - Improve scraping service if needed

### Phase three
  Let's do frontend
  - Setup a basic form to submit a request with ASIN input
  - Setup API call and return result in a nice layout
  - Setup register and log in with PostgreSQL or maybe other databases
  - Setup third-party authentications like Google for login
  - Frontend testing
