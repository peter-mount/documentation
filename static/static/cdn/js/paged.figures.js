function createFigures(config) {
  let figuresElements = document.getElementsByClassName(config.spanClassFigures)
  let arrayIndex = []
  let num = 0

  for (let i = 0; i < figuresElements.length; i++) {
    let elem = figuresElements[i]
    let captions = elem.getElementsByTagName("figcaption")
    if (captions.length == 0) {
      captions = ""
    } else {
      captions = captions[0].textContent
    }

    num++
    let id = elem.id
    if (id === '') {
      id = 'book-figure-' + num
      elem.id = id
    }

    arrayIndex.push({
      id: id,
      num: num,
      elem: elem,
      captions: captions,
    })

  }

  let figuresDiv = document.querySelector(config.figuresElement)

  let ul = document.createElement("ul")
  ul.id = "list-figures-generated"
  ul.classList.add("list-index-generated-1col")
  figuresDiv.appendChild(ul)

  for (let i = 0; i < arrayIndex.length; i++) {
    let figure = arrayIndex[i]

    let li = document.createElement("li")
    li.classList.add("list-index-element")
    li.classList.add("list-figref-element")
    li.dataset.figref = figure.id
    li.dataset.fignum = figure.num
    ul.appendChild(li)

    let span = document.createElement("span")
    span.classList.add("index-value")
    span.textContent = figure.captions
    li.appendChild(span)

    let linkPages = document.createElement("span")
    linkPages.classList.add("link-pages")
    li.appendChild(linkPages)

    span = document.createElement("span")
    span.classList.add("link-page")
    span.innerHTML = '<a href="#' + figure.id + '"></a>'
    linkPages.appendChild(span)
  }
}
