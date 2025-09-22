+++
    "Title": "Dark mode difficulties",
    "Date": "24-5-2024",
    "Description": "How preventing eye sores became a headache"
+++
It's no secret, designing this website gave me many grievances
(in the aesthetics department).
However, I would say none were so troublesome as the matter of proper
utilisation of the colour palette for dark mode, and the
sacrificing of certain design elements in the name of readability.

Designers will quickly find out that they need
to change nearly everything to conform to proper dark mode
design principles. Change the hues of all the colours, decide
on which black is too black, which white is blinding white,
what to do with the drop shadows... 

Thus, a question is imposed. Do we just do a plain re-mapping
of colours, 1:1 so to speak, or do we further analyse the need
for dark mode, and see how to best approach it and it's philosophy.
After all, dark mode is used (in my experience,at least), firstly
for the comfort of the readers' eyes, and secondly due to its
battery and screen saving properties.

If guided by the second point, we should strive for the darkest
shades we can achieve, a pure `#000000` black to be precise, so that
we turn off the pixels for users with OLED displays capable of
per-pixel lighting.
However, the first point mandates than we neither have too big
of a contrast between the background and foreground (text), else
we will cause the user to see "ghosts" around the characters,
as if they were glowing.

Because of the conundrum in the paragraph above, **it's recognised
as a good practice to never use pure blacks or pure whites.**

Side note: because of personal taste, which uses the excuse of
blue light being bad for one's eyes, I will gravitate towards
warmer colors and hues which have less of a blue component in them.

 You may have looked at the image of how this website used to
look in dark mode, and asked yourself: _What was wrong with it?_
A fair question.
_Shadows._

I couldn't have drop shadows in dark mode (as are present so
gorgeously in light mode here), due to the fact it's pretty
damn stupid to have shadows in DARK mode, where there
is perhaps NO LIGHT coming and everything is engulfed in a 
shadow already. Not to mention, the challenges of making sure
shadows are dark enough to be noticed, and the base content isn't 
a shade so undecidedly dark it looks like
[middle gray](https://en.wikipedia.org/wiki/Middle_gray),
just to make sure that the shadows pop.

In the end, I decided to just remove the shadows, and "shift" the
colours by one shade, so everything got a bit darker.

Now, to implement a dark/light mode toggle? That might just be
a story for another day... 