window.addEventListener('load', function () {
  let days = document.getElementsByClassName('day')
  for (let i = 0; i < days.length; i++) {
    days[i].addEventListener('click', function (evt) {
      let links = evt.currentTarget.getElementsByTagName('a')
      if (links.length) location.href = links[0].href
    })
  }
})
