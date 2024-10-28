document.addEventListener('mousemove', function(e) {
    const moveX = (e.clientX / window.innerWidth) - 0.5;
    const moveY = (e.clientY / window.innerHeight) - 0.5;
  
    document.getElementById('layer1').style.transform = `translate(${moveX * 25}px, ${moveY * 25}px)`;
    document.getElementById('layer2').style.transform = `translate(${moveX * 50}px, ${moveY * 50}px)`;
    document.getElementById('layer3').style.transform = `translate(${moveX * 50}px, ${moveY * 50}px)`;
    document.getElementById('layerLogo').style.transform = `translate(${moveX * 25}px, ${moveY * 25}px)`;
    document.getElementById('layerLogoblur').style.transform = `translate(${moveX * 25}px, ${moveY * 25}px)`;
  });