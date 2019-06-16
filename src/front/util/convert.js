export const blobToText = (file) => {
    return new Promise((resolve, reject) => {
        let fr = new FileReader();
        fr.onload = () => {
            resolve(fr.result)
        };
        fr.readAsText(file);
    });
};