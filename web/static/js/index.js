let obj = new Vue({
    el: '#app',
    data: {
        isActive: true,
        fz: '25px',
    },
});

setTimeout(() => {
    obj.isActive = false;
    obj.fz = '30px';
}, 1000);
