async function onSubmit(evt) {
  evt.preventDefault()
  let form = new FormData(evt.currentTarget)
  let response = await fetch('/login', { method: 'POST', body: form })
  if (!response.ok) {
    document.getElementById('login-incorrect').style.display = null
    document.getElementById('password').value = ''
    document.getElementById('password').focus()
    return
  }
  location.href = '/calendar'
}
window.addEventListener('load', function () {
  document.getElementById('form').addEventListener('submit', onSubmit)
})
