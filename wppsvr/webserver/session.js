window.addEventListener('load', function () {
  let sdin = document.getElementById('startDate')
  let sd = sdin.value
  let edin = document.getElementById('endDate')
  let ed = edin.value
  sdin.addEventListener('change', function () {
    sd = sdin.value
  })
  edin.addEventListener('change', function () {
    let nv = Date.parse(edin.value)
    if (ed && sd && nv) {
      let sdv = Date.parse(sd)
      let edv = Date.parse(ed)
      let delta = edv - sdv
      sdin.value = new Date(nv - delta).toISOString().substring(0, 10)
    }
    ed = edin.value
  })

  document.querySelectorAll('input[type=checkbox]').forEach(function (dest) {
    if (!dest.name.startsWith('destbbs.')) return
    let down = document.getElementById('downbbs.' + dest.name.substring(8))
    dest.addEventListener('click', function () {
      if (dest.checked) down.checked = false
    })
    down.addEventListener('click', function () {
      if (down.checked) dest.checked = false
    })
  })

  document.getElementById('anyMessage').addEventListener('click', function () {
    document.getElementById('mtypeRow').style.display = null
    document.getElementById('plainSubjectRow').style.display = 'none'
    document.getElementById('plainBodyRow').style.display = 'none'
    document.getElementById('formBodyRow').style.display = 'none'
    document.getElementById('formImageRow').style.display = 'none'
  })
  document.getElementById('plainMessage').addEventListener('click', function () {
    document.getElementById('mtypeRow').style.display = 'none'
    document.getElementById('plainSubjectRow').style.display = null
    document.getElementById('plainBodyRow').style.display = null
    document.getElementById('formBodyRow').style.display = 'none'
    document.getElementById('formImageRow').style.display = 'none'
  })
  document.getElementById('formMessage').addEventListener('click', function () {
    document.getElementById('mtypeRow').style.display = 'none'
    document.getElementById('plainSubjectRow').style.display = 'none'
    document.getElementById('plainBodyRow').style.display = 'none'
    document.getElementById('formBodyRow').style.display = null
    document.getElementById('formImageRow').style.display = null
  })
})
