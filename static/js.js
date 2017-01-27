
const id = String(Math.random());

// make an ajax call, return a string promise
function ajax(url, data) {
  return new Promise((resolve, reject) => {
    const req = new XMLHttpRequest();
    req.addEventListener('load', () => {
      resolve(JSON.parse(req.responseText));
    });
    req.open('POST', url, true);
    data ? req.send(JSON.stringify(data)) : req.send();
  });
}

function onmessage(message) {
  const response = JSON.parse(message.data);
  const div = document.createElement('div');
  div.textContent = response.text;
  document.getElementById('chat').appendChild(div);
}

function onsend() {
  const input = document.getElementById('input');
  const text = input.value;
  input.value = '';
  input.enabled = false;
  ajax('/rpc/send', {text: text}).then(() => input.enabled = true);
}

ajax('/rpc/login', {id: id}).then(response => {
  const channel = new goog.appengine.Channel(response.token);
  channel.open({
    onmessage: onmessage,
    onopen() {}, onerror() {}, onclose() {},
  });
  document.getElementById('send').addEventListener('click', onsend);
  document.getElementById('input')
      .addEventListener('keypress', e => e.keyCode == 13 && onsend());
});

