Scraper planning:
1. [] Cronjob
2. [] Spider interfaces
3. [] Merger/aggregation functions
4. [x] Channel investigation 
5. [] Worker pool investigation
6. [] Define DB
7. [] Diagram of data flow
8. [] Pipeline diagram.

Tricks to implement in the first version:
- Free ip rotation
- Real user agent (with rotation)
- Set other request headers
- Random intervals (cronjob)
- Referrers
- Avoid invisible links
- Inform about Google Cache website

Pipelines:
- sequential (time management) / concurrent

Make macropromt mentioning:
    - static and dynamic spiders (all might have link channel and output channel).
    - Links might want to remember where they came from.
    - Links from static spiders configured in yaml or env?
    - Routers.
    - Mergers (merge by item id, what should the id be? How is it chosen) merge by bool or by amount of messages.
    - How to communicate. What should channels be? Mergers receive information from multiple spiders, potentially.
    - Pipelines: define flows of spiders with routers and a merger. They also write to DB.
    - Pipelines might be sequential (define time spent) or concurrent (worker pools).

Prompt:
I want to make a Go scraping system using colly to scrape webs that have offers in them. Spiders of this service might be static (they alkways check the same link) which usually check the offer page for links and, perhaps, additional info and dynamic (they receive different links), which usually check the page for info to store in DB and, sometimes, links for other dynamic spiders. Therefore all spiders might be retrieving links, output or both but links for dynamic spiders are not fixed. When sending a link from a spider to another, it should be sent the link that this new link was found in in order to potentially include information in the request headers. For static spiders, should links be in .env or yaml config files? ANother element of this service are routers, which based on the information provided by previous spiders (either links, output or both) decide which spider should analyze the received links using some logic implemented inside it tht might vary from one application to another. Finally, the element that aggregates all the information obtained by the spiders are mergers. Given that a merger might wanna aggregate information coming from several spiders, an item sent by a spider should have an id. Mergers should be able to merge item data when they receive data from an item coming from all involved spiders. Take into account that not all items, despite being the same type might be processed by the same amount of spiders due to routers or other similar factors so the merger should know when an item processing is done and how many spiders contributed. How should the communication between all these element be handled. Take into account that some spiders might be faster than others so there might be several items being processed at once. Finally, all these elements form pipelines, which, after aggregating data write to SQL DB. Pipelines can be either sequential (only one item at a time and some random waiting interval between actions can be set to appear real to the website) or concurrent, using worker pools for spiders.
