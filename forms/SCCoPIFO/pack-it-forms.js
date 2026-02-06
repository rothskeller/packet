/* Copyright 2014 Keith Amidon
   Copyright 2014 Peter Amidon
   Copyright 2018 John Kristian
   Copyright 2025 Steve Roth

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License. */

var required_groups = []
var the_form

var standardAttributes = {
  "cardinal-number": { pattern: "[0-9]*" },
  date: {
    pattern:
      "(0[1-9]|1[012])/(0[1-9]|1[0-9]|2[0-9]|3[01])/[1-2][0-9][0-9][0-9]",
    placeholder: "mm/dd/yyyy",
    cleanupHandler: evt => {
      let value = evt.target.value.replaceAll('-', '/')
      if (/^\d\//.test(value)) value = '0' + value
      if (/^\d\d\/\d\//.test(value)) value = value.substring(0, 3) + '0' + value.substring(3)
      if (/^\d\d\/\d\d\/\d\d$/.test(value)) value = value.substring(0, 6) + '20' + value.substring(6)
      evt.target.value = value
    },
  },
  frequency: { pattern: "[0-9]+(\.[0-9]+)?" },
  "frequency-offset": {
    pattern: "[\\-+]?[0-9]*\\.[0-9]+|[\\-+]?[0-9]+|[\\-+]",
  },
  "nonzero-cardinal-number": { pattern: "[1-9][0-9]*" },
  "phone-number": {
    pattern: "[a-zA-Z ]*([+][0-9]+ )?[0-9][0-9 \\-]*([xX][0-9]+)?",
    placeholder: "000-000-0000 x00",
    cleanupHandler: evt => {
      let value = evt.target.value
      const ext = /[xX][0-9]+$/.exec(value)
      if (ext) value = value.substring(0, ext.index)
      const digits = value.replaceAll(/[^0-9]/g, '')
      if (digits.length === 10) value = digits.substring(0, 3) + '-' + digits.substring(3, 6) + '-' + digits.substring(6)
      if (ext) value += ' ' + ext[0]
      evt.target.value = value
    },
  },
  "real-number": { pattern: "[\\-+]?[0-9]*\\.[0-9]+|[\\-+]?[0-9]+" },
  time: {
    pattern: "([01][0-9]|2[0-3]):?[0-5][0-9]|2400|24:00",
    placeholder: "hh:mm",
    cleanupHandler: evt => {
      let value = evt.target.value
      if (/^\d:?\d\d$/.test(value)) value = '0' + value
      if (/^\d\d\d\d/.test(value)) value = value.substring(0, 2) + ':' + value.substring(2)
      evt.target.value = value
    }
  },
  "zip-code": { pattern: "\\d{5}(?:-\\d{4})?" },
}

// UTILITY FUNCTIONS

function array_for_each(array, func) {
  return Array.prototype.forEach.call(array, func)
}

function setControlValue(control, value) {
  if (control.type === "checkbox" || control.type === "radio")
    control.checked = !!value
  else control.value = value
}

function controlValue(control) {
  if (control.type === "checkbox" || control.type === "radio")
    return control.checked ? "checked" : ""
  else return control.value
}

function anyChildHasValue(elm) {
  let found = false
  elm
    .querySelectorAll(
      "input[type=checkbox],input[type=radio],input[type=text],input:not([type]),select,textarea",
    )
    .forEach((control) => {
      if (controlValue(control)) found = true
    })
  return found
}

// A Conditional interprets the conditional in an element attribute
// (hidden-until, required-if, allowed-if) and dispatches 'change' events when
// the state of the conditional changes.
class Conditional extends EventTarget {
  static conditionals = {};
  static getOrMake(elm, attr) {
    const cstr = elm.getAttribute(attr)
    if (!cstr) return null
    return this.conditionals[cstr] || new Conditional(elm, attr)
  }
  constructor(elm, attr) {
    super()
    const form = elm.closest("form")
    if (!form) throw `${attr} outside of a form`
    const cstr = elm.getAttribute(attr)
    if (!cstr) return null
    const parts = cstr.split("=", 2)
    this.field = form[parts[0]]
    if (!this.field)
      throw `${attr}="${cstr}": no such form element "${parts[0]}"`
    if (this.field.type === "checkbox") this.test = () => this.field.checked
    else if (parts.length > 1) this.test = () => this.field.value == parts[1]
    else this.test = () => !!this.field.value
    this.was = this.test()
    if (this.field instanceof RadioNodeList)
      this.field.forEach((b) => {
        b.addEventListener("change", this.onChange.bind(this))
      })
    else if (this.field.type === "checkbox" || this.field.type === "radio")
      this.field.addEventListener("change", this.onChange.bind(this))
    else this.field.addEventListener("input", this.onChange.bind(this))
    this.constructor.conditionals[cstr] = this
  }
  onChange() {
    const is = this.test()
    if (is != this.was) {
      this.was = is
      this.dispatchEvent(new CustomEvent("change", { detail: is }))
    }
  }
}

// SHARED FUNCTIONS (idempotent, called in startup and event handlers)

// Enables or disables the submit buttons based on form validity.
function adjust_submit() {
  var valid = the_form.checkValidity()
  document.querySelector("#button-header").classList.toggle("valid", valid)
  document.querySelector("#submit").disabled = !valid
  var invalid_example = document.querySelector("#invalid-example")
  invalid_example.hidden = valid
  invalid_example.classList.toggle("hidden", valid)
  return valid
}

// Adjusts the pattern of a text field based on its type and required state.
// For a textarea, it adjusts the invalid class.
function adjust_pattern(input) {
  var pattern = input.pattern
  for (var s in standardAttributes) {
    if (input.classList.contains(s)) {
      var standard = standardAttributes[s]
      if (standard) pattern = standard.pattern
    }
  }
  if (input.type == "textarea") {
    input.classList.toggle(
      "invalid",
      input.required && (!input.value || !input.value.trim()),
    )
  } else if (input.type == "text") {
    if (pattern == "\\s*\\S.*") pattern = ""
    if (input.required) {
      if (!pattern) {
        pattern = "\\s*\\S.*" // not all white space
      }
    } else if (pattern) {
      pattern += "|\\s*" // all white space
    }
  }
  if (pattern) {
    if (input.classList.contains("clearable")) {
      pattern += "|\\{CLEAR\\}"
    }
    if (input.pattern != pattern) {
      input.pattern = pattern
    }
  } else if (input.pattern) {
    input.removeAttribute("pattern")
  }
}

// Adjusts the required flags on all controls in a required group.
function adjust_required_group(group) {
  if (group.querySelector(":checked") || group.closest("[hidden]") || !group.classList.contains('required-group')) {
    group.querySelectorAll("input[type=checkbox]:required").forEach((r) => {
      r.required = false
    })
    group.classList.remove("invalid")
  } else {
    let haveRequired = false,
      seenRFC = false
    group
      .querySelectorAll('input[type="checkbox"],input[type="radio"]')
      .forEach((r) => {
        if (r.disabled) return
        r.required = true
        haveRequired = true
      })
    group.classList.toggle("invalid", haveRequired)
  }
}

// Adjusts the required flags on all controls in required groups.
function adjust_required_groups() {
  required_groups.forEach(adjust_required_group)
  adjust_submit()
}

// Adjusts the required and disabled flags of controls based on required-if,
// allowed-if, and else-disallowed attributes.
function adjust_required_disabled() {
  document
    .querySelectorAll(
      "input[type=checkbox],input[type=radio],input[type=text],input:not([type]),select,textarea",
    )
    .forEach((elm) => {
      if (elm.closest("[hidden]")) return
      let state
      if (elm.hasAttribute("required-if")) {
        if (elm.hasAttribute("allowed-if"))
          throw "control cannot have both required-if and allowed-if attributes"
        const cond = Conditional.getOrMake(elm, "required-if")
        if (cond.test()) state = "required"
        else if (elm.hasAttribute("else-disallowed")) state = "disallowed"
        else state = "optional"
      } else if (elm.hasAttribute("allowed-if")) {
        const cond = Conditional.getOrMake(elm, "allowed-if")
        state = cond.test() ? "optional" : "disallowed"
      }
      switch (state) {
        case "required":
          elm.required = true
          elm.disabled = false
          if (elm.hasAttribute("disallowed-value")) {
            setControlValue(elm, elm.getAttribute("disallowed-value"))
            elm.removeAttribute("disallowed-value")
          }
          break
        case "optional":
          elm.required = false
          elm.disabled = false
          if (elm.hasAttribute("disallowed-value")) {
            setControlValue(elm, elm.getAttribute("disallowed-value"))
            elm.removeAttribute("disallowed-value")
          }
          break
        case "disallowed":
          elm.required = false
          if (elm.type === "checkbox" || elm.type === "radio") {
            if (elm.checked) elm.setAttribute("disallowed-value", "checked")
            elm.checked = false
          } else if (elm.value) {
            elm.setAttribute("disallowed-value", elm.value)
            elm.value = ""
          }
          elm.disabled = true
          break
      }
      adjust_pattern(elm)
    })
  document
    .querySelectorAll(".required-group,.was-required-group")
    .forEach((elm) => {
      if (elm.closest("[hidden]")) return
      if (!elm.hasAttribute("required-if")) return
      const cond = Conditional.getOrMake(elm, "required-if")
      if (cond.test()) {
        elm.classList.remove('was-required-group')
        elm.classList.add('required-group')
      } else {
        elm.classList.remove('required-group')
        elm.classList.add('was-required-group')
      }
    })
  adjust_required_groups()
}

// EVENT HANDLERS

// Triggered on any form input.
function on_form_input() {
  adjust_submit()
}

// Triggered when a hidden-until element's condition becomes true.
function on_hidden_until(evt, elm) {
  if (!evt.target.test()) return
  elm.removeAttribute("hidden")
  elm.querySelectorAll("[hidden-save-required]").forEach((control) => {
    control.setAttribute(
      "required",
      control.getAttribute("hidden-save-required"),
    )
    control.removeAttribute("hidden-save-required")
  })
  adjust_required_disabled()
}

// Triggered when a control in a required-group changes.
function on_required_group_change(evt) {
  adjust_required_group(evt.target.closest(".required-group"))
  adjust_submit()
}

// Triggered by Reset Form button.
function on_reset() {
  the_form.reset()
  reset_form()
}

// Triggered by Submit or Save button.
async function on_submit(evt) {
  const submit = evt.target.id == "submit"
  if (submit && !adjust_submit()) return
  const fd = new FormData(the_form)
  if (submit) fd.set("readyToSend", "true")
  const resp = await fetch(the_form.action, {
    method: "POST",
    body: fd,
    redirect: "manual",
  })
  if (resp.status == 204) {
    const action = resp.headers.get("X-Packet-Action")
    if (action.startsWith("redirect:")) location.href = action.substring(9)
    else {
      window.opener.childAction(resp.headers.get("X-Packet-Action"))
      window.close()
    }
  } else document.getElementById("error").textContent = await resp.text()
}

// Triggered by Show PDF button.
function on_show_pdf() {
  window.open(the_form.dataset.pdfUrl, "_blank")
}

// RESET FUNCTIONS (run at startup and whenever the form is reset)

// Re-hides elements with hidden-until attributes that are not satisfied.
function reset_hidden_until() {
  document.querySelectorAll("[hidden-until]").forEach((elm) => {
    const cond = Conditional.getOrMake(elm, "hidden-until")
    if (cond.test() || anyChildHasValue(elm)) return
    elm.setAttribute("hidden", "")
    elm.querySelectorAll("[required]").forEach((control) => {
      control.setAttribute(
        "hidden-save-required",
        control.getAttribute("required"),
      )
      control.removeAttribute("required")
    })
  })
}

function reset_form() {
  reset_hidden_until()
  adjust_required_disabled()
}

// SETUP FUNCTIONS (run only once)

// Ping the server periodically, to keep it alive while this page is open.
function start_pings() {
  setInterval(function ping() {
    var img = new Image()
    // To discourage caching, use a new query string for each ping.
    img.src = "/ping?i=" + Math.random()
    img = undefined
  }, 30000) // every 30 seconds
}

// Sets up listeners on elements with unsatisfied hidden-until attributes.
function setup_hidden_until() {
  document.querySelectorAll("[hidden-until]").forEach((elm) => {
    const cond = Conditional.getOrMake(elm, "hidden-until")
    if (cond.test() || anyChildHasValue(elm)) return
    cond.addEventListener("change", (evt) => {
      on_hidden_until(evt, elm)
    })
  })
}

// Returns whether an element has required-if or allowed-if (or, implicitly,
// else-disallowed) attributes.
function has_conditionals(elm) {
  return elm.hasAttribute("required-if") || elm.hasAttribute("allowed-if")
}

// Copies the required-if, allowed-if, and/or else-disallowed attributes from
// the nearest ancestor that has them, if any.
function inherit_conditionals(elm) {
  if (has_conditionals(elm)) return
  for (let p = elm.parentElement; p; p = p.parentElement) {
    if (has_conditionals(p)) {
      if (p.hasAttribute("required-if"))
        elm.setAttribute("required-if", p.getAttribute("required-if"))
      if (p.hasAttribute("allowed-if"))
        elm.setAttribute("allowed-if", p.getAttribute("allowed-if"))
      if (p.hasAttribute("else-disallowed"))
        elm.setAttribute("else-disallowed", p.getAttribute("else-disallowed"))
      return
    }
  }
}

// Sets up conditional presence requirements based on the required-if,
// allowed-if, and else-disallowed attributes.
function setup_required_if() {
  document
    .querySelectorAll("input[type=checkbox],input[type=radio]")
    .forEach((elm) => {
      inherit_conditionals(elm)
      elm.addEventListener("change", adjust_required_disabled)
    })
  document
    .querySelectorAll("input[type=text],input:not([type]),select,textarea")
    .forEach((elm) => {
      inherit_conditionals(elm)
      elm.addEventListener("input", adjust_required_disabled)
    })
  document.querySelectorAll('.required-group').forEach(elm => {
    inherit_conditionals(elm)
  })
  adjust_required_disabled()
}

// Sets up required groups.
function setup_required_groups() {
  required_groups = Array.from(document.querySelectorAll(".required-group,.was-required-group"))
  required_groups.forEach((g) => {
    g.addEventListener("change", on_required_group_change)
  })
}

// Sets the properties of an input that only need to be set once on startup.
// These include placeholder, title, shift-click handler for radio buttons.
function setup_input_once(input) {
  if (!input.placeholder) {
    for (var s in standardAttributes) {
      if (input.classList.contains(s)) {
        var placeholder = standardAttributes[s]?.placeholder
        if (placeholder) input.placeholder = placeholder
        var cleanup = standardAttributes[s]?.cleanupHandler
        if (cleanup) input.addEventListener('change', evt => {
          cleanup(evt)
          adjust_submit()
        })
      }
    }
  }
  if (!input.title && input.placeholder) {
    input.title = input.placeholder
  }
  if (input.type == "radio") {
    input.addEventListener("click", (evt) => {
      if (evt.shiftKey && !input.required) input.checked = false
    })
  }
}

function setup_buttons() {
  document.getElementById("submit").addEventListener("click", on_submit)
  document.getElementById("reset").addEventListener("click", on_reset)
  const save = document.getElementById("save")
  if (save) save.addEventListener("click", on_submit)
}

function set_initial_focus() {
  const invalid = document.querySelector('input:invalid,select:invalid,textarea:invalid')
  if (invalid) invalid.focus()
}

window.addEventListener("load", function () {
  the_form = document.getElementById("the-form")
  start_pings()
  setup_hidden_until()
  setup_required_if()
  setup_required_groups()
  array_for_each(the_form.elements, setup_input_once)
  setup_buttons()
  the_form.addEventListener("input", on_form_input)
  reset_form()
  set_initial_focus()
})
