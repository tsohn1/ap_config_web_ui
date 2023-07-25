let dropdownItems = document.querySelectorAll(".dropdown-item");
let currentLang = Cookies.get("language");
let korean = document.querySelector("[data-language = 'kr']");
let english = document.querySelector("[data-language = 'en']");
if (currentLang === undefined) {
  currentLang = "en";
}
let operation = document.querySelector("#operationNav")
let network = document.querySelector("#networkNav")
let lang = document.querySelector("#languageDropdown")
console.log(operation, network, lang)
if (currentLang === "kr") {
  dropdownItems.forEach(item => item.classList.remove("active"));
  korean.classList.add("active");
  operation.textContent = "오퍼레이션"
  network.textContent = "네트워크"
  lang.textContent = "언어"
} else {
  dropdownItems.forEach(item => item.classList.remove("active"));
  english.classList.add("active");
}

