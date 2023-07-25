let dropdownItems = document.querySelectorAll(".dropdown-item");
let currentLang = Cookies.get("language");
let korean = document.querySelector("[data-language = 'kr']");
let english = document.querySelector("[data-language = 'en']");
if (currentLang === undefined) {
  currentLang = "en";
}
console.log(currentLang);
console.log(dropdownItems);
console.log(korean, english);
if (currentLang === "kr") {
  dropdownItems.forEach(item => item.classList.remove("active"));
  korean.classList.add("active");
} else {
  dropdownItems.forEach(item => item.classList.remove("active"));
  english.classList.add("active");
}
