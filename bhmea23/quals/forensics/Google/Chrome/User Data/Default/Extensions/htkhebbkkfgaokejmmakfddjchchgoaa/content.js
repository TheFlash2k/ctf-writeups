const userContent = document.body.innerText;

chrome.runtime.onMessage.addListener(function (message, sender, sendResponse) {
  if (message.action === "getContent") {
    sendResponse({ content: userContent });
  }
});