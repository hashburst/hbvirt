// static/js/scripts.js
document.getElementById('registration-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const email = document.getElementById('email').value;
    const apikey = document.getElementById('apikey').value;
    const referralCode = document.getElementById('referral_code').value;
    const dogeWallet = document.getElementById('doge_wallet').value;
    const btcWallet = document.getElementById('btc_wallet').value;

    const response = await fetch('/verify', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            email: email,
            apikey: apikey,
            referral_code: referralCode,
            doge_wallet: dogeWallet,
            btc_wallet: btcWallet
        })
    });

    const result = await response.json();
    if (result.success) {
        document.getElementById('status-message').innerText = 'Licenza verificata! Puoi usare il programma gratuitamente per 6 mesi.';
    } else {
        document.getElementById('status-message').innerText = 'Verifica fallita: ' + result.error;
    }
});
