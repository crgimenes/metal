# metal
Bare metal software


Video colors 

CGA	EGA	 VGA	    RGB	         Web    	Example
0x0	0x0	 0,0,0	    0,0,0	     #000000	black
0x1	0x1	 0,0,42	    0,0,170      #0000aa	blue
0x2	0x2	 00,42,00	0,170,0	     #00aa00	green
0x3	0x3	 00,42,42	0,170,170    #00aaaa	cyan
0x4	0x4	 42,00,00	170,0,0	     #aa0000	red
0x5	0x5	 42,00,42	170,0,170    #aa00aa	magenta
0x6	0x14 42,21,00	170,85,0     #aa5500	brown
0x7	0x7	 42,42,42	170,170,170	 #aaaaaa	gray
0x8	0x38 21,21,21	85,85,85	 #555555	dark gray
0x9	0x39 21,21,63	85,85,255	 #5555ff	bright blue
0xA	0x3A 21,63,21	85,255,85	 #55ff55	bright green
0xB	0x3B 21,63,63	85,255,255	 #55ffff	bright cyan
0xC	0x3C 63,21,21	255,85,85	 #ff5555	bright red
0xD	0X3D 63,21,63	255,85,255	 #ff55ff	bright magenta
0xE	0x3E 63,63,21	255,255,85	 #ffff55	Yellow
0xF	0x3F 63,63,63	255,255,255	 #ffffff	white

Video memory

+---------+----------+
| 8 bits  | 8 bits   | 4 bits to front color and 4 bits to back color
+---------+----------+
|char code|color code| 
+---------+----------+
