export const joinMsg =  (token, name) => {
    return {
        'Join': {
            'Token': token,
            'Name': name,
        }
    }
};

export const turnEndMsg =  (token, board) => {
    return {
        'EndTurn': {
            'Token': token,
            'Board': board,
        }
    }
};