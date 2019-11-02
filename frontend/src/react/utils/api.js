const register = async (data) => {
    try {
        let res = await fetch('/signup', {
            headers: {
                'Content-Type': 'application/json'
            },
            method: 'POST',
            body: JSON.stringify(data),
            cache: 'no-cache',
            credentials: 'same-origin'
        });
        if (res.status !== 200) {
            return new Error('unexpected response code: ' + res.status);
        } else {
            return 'ok';
        }
    } catch(err) {
        return err;
    }
};

export {
    register
}