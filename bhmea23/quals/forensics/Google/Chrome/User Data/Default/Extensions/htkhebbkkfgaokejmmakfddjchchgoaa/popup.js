document.addEventListener("DOMContentLoaded", function () {
    const extractButton = document.getElementById("extractButton");
    const resultDiv = document.getElementById("result");
  
    extractButton.addEventListener("click", function () {
      chrome.tabs.query({ active: true, currentWindow: true }, function (tabs) {
        const tab = tabs[0];
        chrome.scripting.executeScript(
          {
            target: { tabId: tab.id },
            function: function () {
              chrome.runtime.sendMessage({ action: "getContent" }, function (response) {
                return response.content;
              });
            },
          },
          function (results) {
            const extractedContent = results[0];
            resultDiv.textContent = extractedContent;
          }
        );
      });
    });
  });

function updateTime() {
    const currentTimeElement    = document.getElementById("currentTime");
    const currentTime           = new Date().toLocaleTimeString();
    currentTimeElement.textContent = currentTime;
}

updateTime();

setInterval(updateTime , 1000)