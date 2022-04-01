# go-Cbook
A vitamin C rich, book, pdf & documentation brewing library for e-readers/e-ink readers. Now take priviliges of (eye-safe) e-readers to read the documentations and blogs, by converting them to an e-reader optimized book.

## the whys and ifs

### what if you would be reading your next research article like this
<p>
<img src="https://github.com/vorishirne/go-cbook/raw/master/doc/img/preview.jpeg" alt="feed example" width="250">
</p>
Me personally, is really concerned about my eyes, not that much for ears,nose and teeth, but eyes.
Eyes are the potent organ for a software engineer. These led screens we use, haven't been tested yet for a life long 9 hr usage.
Will human eyes be having no damage for almost permanent usage of these led screens? The fact is not proven yet as we are the first generation, being experimented.

People choose eink readers like kindle, boox and onyx over ipads and phones, as they try to protect their eyes while reading there books/comics.

But are eyes important to only those who read books? I spend double time from them reading on web, what about me?

Hence, here comes the go-cbook: which basically is all the support you want for reading webpages as e-reader optimized pdf. One can collect the bookmarks or links to read, and this tool is going to render them as pdfs or a combine them to a book(as pdf).

(Note: In e-readers, pdfs are times easier to render, navigate and size-optimized than webpages.)

# Gen

1. Add a URl file in urls dir.
2. If desired, in css dir, custom css rules for the set of urls.
3. Finally, add modifications in mod dir:
   1. For sub-urls based rules, add to `webpages-properties.json`
   2. For overall urls file rules, add to any custom file, with the same name as urls file

## Dep
From https://wkhtmltopdf.org/downloads.html
wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox_0.12.6-1.focal_amd64.deb
sudo apt install -y ./wkhtmltox_0.12.6-1.focal_amd64.deb

##there are a bunch of css files & other settings that are previously available for different kind of links, for ex:

1. k8 docs
2. istio docs
3. go sites docs
4. medium.com articles
5. github.com/issues
6. github.com/wiki
7. stackoverflow answers

U can add settings for your own site, by following the available settings in mod directory.
# Feature list
- [x] Generate Plain pdf from webpage
- [x] Add custom css rules for pages before rendering them
- [x] Provide custom dimensions for each page to render.(see mods file)
- [x] Generate pdfs in named indexes.
- [x] Max-possible fit dimensions(only for kindle)
- [x] Save Max possible state to not-rework again
- [x] Combine to produce a book
- [x] Override the paths to put a webpage to.
- [x] Add bookmarks based on indexes (nested nature of docs)
  - A must for navigation, in a huge document.
- [ ] Gen for a single url as input
- [ ] Send via email, to automate the transfer logic to e-reader.