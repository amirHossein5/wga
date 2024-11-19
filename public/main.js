const ws = new WebSocket('ws://127.0.0.1:8080/update-player');

const connecting = document.querySelector('#connecting');
const playgroundCanvas = document.querySelector('#playground');
const playgroundContext = playgroundCanvas.getContext('2d');
const controlButtons = document.querySelector('#control-buttons');

ws.addEventListener('message', function (event) {
    const data = JSON.parse(event.data);

    clearCanvas();

    assert(data instanceof Object, 'data is not an object', data);
    assert('players' in data, 'failed to receieve players', data);
    assert('current_player' in data, 'failed to receieve current_player', data);

    data.players.forEach(function (player) {
        drawPlayer(player, data['current_player'].id === player.id);
    });

    connecting.style.display = 'none';
    playgroundCanvas.style.display = 'block';
    controlButtons.style.display = 'flex';
});

function drawPlayer(player, isCurrentPlayer) {
    playgroundContext.fillStyle = player['color'];

    if (isCurrentPlayer) {
        // playgroundContext.strokeStyle = 'red';
        // playgroundContext.lineWidth = '2';
        // playgroundContext.rect(player['position_x'], player['position_y'], 40, 40);
    }

    playgroundContext.fillRect(player['position_x'], player['position_y'], 40, 40);

    // playgroundContext.fill();
    // if (isCurrentPlayer) playgroundContext.stroke();
}

function clearCanvas() {
    playgroundContext.clearRect(0, 0, playgroundCanvas.width, playgroundCanvas.height);
}

function assert(thurhy, msg, ...data) {
    if (!thurhy) {
        console.error(msg, ...data);
        throw new Error(msg);
    }
}

document.querySelector('button.to-right').addEventListener('click', moveToRight);
document.querySelector('button.to-left').addEventListener('click', moveToLeft);
document.querySelector('button.to-top').addEventListener('click', moveToTop);
document.querySelector('button.to-down').addEventListener('click', moveToDown);

document.onkeydown = function(e) {
    if (e.key === 'ArrowRight') {
        moveToRight();
    } else if (e.key === 'ArrowLeft') {
        moveToLeft();
    } else if (e.key === 'ArrowUp') {
        moveToTop();
    } else if (e.key === 'ArrowDown') {
        moveToDown();
    }
}

function moveToRight() {
    ws.send('to-right');
}

function moveToLeft() {
    ws.send('to-left');
}

function moveToTop() {
    ws.send('to-top');
}

function moveToDown() {
    ws.send('to-down');
}