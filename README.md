# screen-to-pdf

Scripts for automating converting slides, ebooks or any other web format to PDF when only visualization is possible.

The main script is in charge of taking a screenshot of an are of the screen, doing something that the user must define to go to the next slide, page, etc and stop when the last image is exactly the same as the previous one. After that all screenshot are converted to PDF.  

The variables for configuring the area to screenshot, the margins of the PDF or the behavior to go to the next slide are at the start of the script.  

In order to make the script work there is a need to create the temporary folders for storing the screenshots and the final PDF if they are not created. The path for those folders in in the configurable variables of the script.

Another script in the folder `getmouselocation` track and print the location of the mouse pointer to help with the setting of the main script variables.
