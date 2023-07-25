document.addEventListener("DOMContentLoaded", function () {
  const dropdownItems = document.querySelectorAll(".dropdown-item");

  dropdownItems.forEach(item => {
    item.addEventListener("click", function (event) {
      event.preventDefault();
      selectedLanguage = this.getAttribute("data-language");
      dropdownItems.forEach(item => item.classList.remove("active"));
      this.classList.add("active");
      Cookies.set("language", selectedLanguage, { expires: 30});
      console.log(selectedLanguage)
      window.location.reload();
    });
  });
});
