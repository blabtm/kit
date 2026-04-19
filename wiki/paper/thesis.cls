\NeedsTeXFormat{LaTeX2e}
\ProvidesClass{thesis}[2025/06/06 Bachelor's Thesis]

\LoadClass[
  fontsize=14pt,
  titlepage,
  usegeometry,
  toc=flat,
  chapterentrydots=true,
  captions=tableheading,
]{scrreprt}

\RequirePackage[
  a4paper,
  nomarginpar,
  noheadfoot,
  hmargin={3cm, 1cm},
  vmargin=2cm,
  footskip=1cm,
]{geometry}

\RequirePackage{fontspec}
\setmainfont{Liberation Serif}
\setsansfont{Liberation Sans}
\setmonofont{Liberation Mono}

\RequirePackage{polyglossia}
\setdefaultlanguage{russian}
\setotherlanguage{english}

\RequirePackage{hyperref}
\RequirePackage{indentfirst}
\RequirePackage{pdfpages}
\RequirePackage{lastpage}

\RequirePackage{totcount}
\regtotcounter{chapter}
\regtotcounter{figure}

\renewcommand\raggedsection{\centering}

\setkomafont{chapter}{\fontsize{18}{18}\selectfont\rmfamily\MakeUppercase}
\setkomafont{section}{\fontsize{16}{16}\selectfont\rmfamily\MakeUppercase}
\setkomafont{subsection}{\fontsize{16}{16}\selectfont\rmfamily\MakeUppercase}
\setkomafont{subsubsection}{\fontsize{14}{14}\selectfont\rmfamily\MakeUppercase}
\setkomafont{paragraph}{\fontsize{14}{14}\selectfont\rmfamily}

%\renewcommand*{\chapterformat}{\thechapter\enskip}
%\renewcommand*{\sectionformat}{\thesection\enskip}
%\renewcommand*{\subsectionformat}{\thesubsection\enskip}
%\renewcommand*{\subsubsectionformat}{\thesubsubsection\enskip}

\setkomafont{chapterentry}{\rmfamily}

\renewcommand{\addchaptertocentry}[2]{%
\IfArgIsEmpty{#1}
 {\addcontentsline{toc}{chapter}{\protect\nonumberline\texorpdfstring{\protect\MakeUppercase{#2}}{#2}}}
 {\addcontentsline{toc}{chapter}{\protect\numberline{#1}\texorpdfstring{\protect\MakeUppercase{#2}}{#2}}}%
}

\renewcommand{\addsectiontocentry}[2]{%
\IfArgIsEmpty{#1}
 {\addcontentsline{toc}{section}{\protect\nonumberline\texorpdfstring{\protect\MakeUppercase{#2}}{#2}}}
 {\addcontentsline{toc}{section}{\protect\numberline{#1}\texorpdfstring{\protect\MakeUppercase{#2}}{#2}}}%
}

\RedeclareSectionCommand[
  beforeskip=0pt,
  afterskip=18pt
]{chapter}

\RedeclareSectionCommand[
  beforeskip=18pt,
  afterskip=12pt
]{section}

\RedeclareSectionCommand[
  beforeskip=12pt,
  afterskip=6pt
]{subsection}

\RedeclareSectionCommand[
  beforeskip=12pt,
  afterskip=6pt
]{subsubsection}

\setcounter{secnumdepth}{3}
\setcounter{tocdepth}{1}

\counterwithout{figure}{chapter}
\counterwithout{figure}{section}
\counterwithout{table}{chapter}
\counterwithout{table}{section}
\counterwithout{equation}{chapter}
\counterwithout{equation}{section}

\RequirePackage{setspace}
\onehalfspacing

\RequirePackage{enumitem}
\setlist[itemize]{noitemsep, topsep=6pt}
\setlist[enumerate]{noitemsep, topsep=6pt}

\RequirePackage{minted}
\setminted{fontsize=\fontsize{10}{10}\selectfont\linespread{1.2}}
\setminted{obeytabs=true}
\setminted{tabsize=4}
\setminted{autogobble}

\RequirePackage{tabularx}
\newcolumntype{Y}{>{\centering\arraybackslash}X}

\RequirePackage{amsmath}
\RequirePackage{syntax}
\RequirePackage{graphicx}
\RequirePackage{float}

\renewcommand*{\figureformat}{Рисунок~\thefigure}
\renewcommand*{\tableformat}{Таблица~\thetable}

\setlength{\intextsep}{12pt}
\setlength{\parindent}{1cm}

\setlength{\belowdisplayskip}{0pt} \setlength{\belowdisplayshortskip}{0pt}
\setlength{\abovedisplayskip}{0pt} \setlength{\abovedisplayshortskip}{0pt}

\RequirePackage{caption}
\captionsetup[figure]{
    skip=6pt,
    labelsep=endash
}
\captionsetup[table]{
    skip=6pt, 
    position=above, 
    justification=raggedright,
    singlelinecheck=false,
    labelsep=endash
}

\let\oldalign\align
\let\oldendalign\endalign
\def\align{\vspace{-12pt}\begingroup\oldalign}
\def\endalign{\oldendalign \endgroup}

\RequirePackage[
  backend=biber,
  bibencoding=utf8,
  bibstyle=gost-numeric,
]{biblatex}