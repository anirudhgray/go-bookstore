// let playerToggle = document.getElementById('toggleLottie');

// playerToggle.addEventListener('ready', () => {
//   LottieInteractivity.create({
//     player: '#toggleLottie',
//     mode: 'cursor',
//     actions: [
//       {
//         type: 'toggle',
//       },
//     ],
//   });
// });

LottieInteractivity.create({
  player: '#toggleLottie',
  mode: 'cursor',
  actions: [
    {
      type: 'toggle',
      frames: [0, 60],
    },
  ],
});

document.getElementById('toggleLottieContainer').onclick = () => {
  if (localStorage.getItem('dark') === '"true"') {
    localStorage.setItem('dark', '"false"');
    document.querySelector('body').classList.remove('dark');
  } else {
    localStorage.setItem('dark', '"true"');
    document.querySelector('body').classList.add('dark');
  }
};

if (localStorage.getItem('dark') === '"true"') {
  // console.log('dark');
  document.querySelector('body').classList.add('dark');
} else {
  // console.log('light');
  document.querySelector('body').classList.remove('dark');
}
