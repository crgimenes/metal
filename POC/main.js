window.addEventListener("keydown", function (event) {

  if (event.defaultPrevented) {
    return; // Should do nothing if the default action has been cancelled
  }

  var handled = false;
  if (event.key !== undefined) {
    console.log(event.key)
    // Handle the event with KeyboardEvent.key and set handled true.
  } 
  if (event.keyIdentifier !== undefined) {
    console.log(event.keyIdentifier)

    // Handle the event with KeyboardEvent.keyIdentifier and set handled true.
  } 
  if (event.keyCode !== undefined) {
    console.log(event.keyCode)

    // Handle the event with KeyboardEvent.keyCode and set handled true.
  }

  if (handled) {
    // Suppress "double action" if event handled
    event.preventDefault();
  }
}, true);

window.onload = function() {
    // Get the canvas and context
    var canvas = document.getElementById("viewport"); 
    var context = canvas.getContext("2d");
 
    // Define the image dimensions
    var width = canvas.width;
    var height = canvas.height;
 
    // Create an ImageData object
    var imagedata = context.createImageData(width, height);
 
    // Create the image
    function createImage(offset) {
        // Loop over all of the pixels
        for (var x=0; x<width; x++) {
            for (var y=0; y<height; y++) {
                // Get the pixel index
                var pixelindex = (y * width + x) * 4;
 
                // Generate a xor pattern with some random noise
                var red = ((x+offset) % 256) ^ ((y+offset) % 256);
                var green = ((2*x+offset) % 256) ^ ((2*y+offset) % 256);
                var blue = 50 + Math.floor(Math.random()*100);
 
                // Rotate the colors
                blue = (blue + offset) % 256;
 
                // Set the pixel data
                imagedata.data[pixelindex] = red;     // Red
                imagedata.data[pixelindex+1] = green; // Green
                imagedata.data[pixelindex+2] = blue;  // Blue
                imagedata.data[pixelindex+3] = 255;   // Alpha
            }
        }
    }


    // Main loop
    function main(tframe) {
        // Request animation frames
        window.requestAnimationFrame(main);
 
        // Create the image
        createImage(Math.floor(tframe / 10));
 
        // Draw the image data to the canvas
        context.putImageData(imagedata, 0, 0);
    }
 
    // Call the main loop
    main(0);
};