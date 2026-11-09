document.addEventListener("DOMContentLoaded", () => {
  const tabs = document.querySelectorAll("[data-tab]");
  const contents = document.querySelectorAll("[data-tab-content]");

  tabs.forEach((tab) => {
    tab.addEventListener("click", () => {
      const target = tab.dataset.tab;

      tabs.forEach((item) => {
        item.classList.remove("active");
      });

      contents.forEach((content) => {
        content.classList.remove("active");
      });

      tab.classList.add("active");

      document
        .querySelector(`[data-tab-content="${target}"]`)
        ?.classList.add("active");
    });
  });
});

document.addEventListener("DOMContentLoaded", () => {
  const tabs = document.querySelectorAll("[data-tab]");
  const contents = document.querySelectorAll("[data-tab-content]");

  tabs.forEach((tab) => {
    tab.addEventListener("click", () => {
      const target = tab.dataset.tab;

      const targetContent = document.querySelector(
        `[data-tab-content="${target}"]`,
      );

      if (!targetContent) return;

      tabs.forEach((item) => {
        item.classList.toggle("active", item === tab);
      });

      contents.forEach((content) => {
        content.classList.toggle("active", content === targetContent);
      });
    });
  });
});
