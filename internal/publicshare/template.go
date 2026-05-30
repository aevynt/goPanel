package publicshare

import (
	"fmt"
	"html"
	"strings"
)

func GalleryPage(share *Share, files []FileEntry, baseURL string) string {
	title := html.EscapeString(share.Title)
	desc := html.EscapeString(share.Description)

		var descHTML string
	if desc != "" {
		descHTML = `<p>` + desc + `</p>`
	}

	var thumbItems []string
	for _, f := range files {
		if f.IsDir {
			continue
		}
		imgPath := "/p/" + share.Folder + "/" + f.Name
		thumbSrc := imgPath
		if f.Thumbnail != "" {
			thumbSrc += "?thumb=1"
		}
		name := html.EscapeString(f.Name)
		thumbItems = append(thumbItems, fmt.Sprintf(`<div class="thumb-item" data-src="%s">
  <div class="thumb-wrap"><div class="thumb-placeholder"></div><img data-src="%s" alt="%s"></div>
  <div class="thumb-name">%s</div>
</div>`, imgPath, thumbSrc, name, name))
	}

	gallery := strings.Join(thumbItems, "\n      ")

	return `<!DOCTYPE html>
<html lang="vi">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=no">
<title>` + title + `</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#0d0d0c;color:#e8e6dc;min-height:100vh;display:flex;flex-direction:column}
.header{padding:32px 24px 16px;max-width:1280px;margin:0 auto;width:100%}
.header h1{font-size:1.75rem;font-weight:600;letter-spacing:-0.02em;margin-bottom:4px}
.header p{color:#9a9990;font-size:14px}
.grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(240px,1fr));gap:12px;padding:0 24px;max-width:1280px;margin:0 auto;width:100%;flex:1;align-content:start;align-items:start}
.thumb-item{cursor:pointer;border-radius:10px;overflow:hidden;background:#1a1a18;transition:transform .15s ease,box-shadow .15s ease}
.thumb-item:hover{transform:translateY(-2px);box-shadow:0 8px 24px rgba(0,0,0,.4)}
.thumb-wrap{width:100%;aspect-ratio:4/3;overflow:hidden;background:#1a1a18;position:relative}
.thumb-placeholder{width:100%;height:100%;position:absolute;inset:0;background:#1a1a18}
.thumb-wrap img{width:100%;height:100%;object-fit:cover;display:block;position:relative;z-index:1;opacity:0;transition:opacity .3s ease}
.thumb-wrap img.loaded{opacity:1}
.thumb-name{padding:8px 12px;font-size:12px;color:#9a9990;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}

/* lightbox */
.lightbox{display:none;position:fixed;inset:0;z-index:9999;background:rgba(0,0,0,.94);flex-direction:column}
.lightbox.active{display:flex}
.lightbox-top{position:absolute;top:0;left:0;right:0;display:flex;align-items:center;padding:12px 16px;z-index:10;height:52px}
.lightbox-top .lb-name{flex:1;font-size:14px;color:#e8e6dc;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;margin-right:12px}
.lb-action{background:none;border:none;color:#9a9990;cursor:pointer;font-size:22px;padding:8px;line-height:1;transition:color .15s;-webkit-tap-highlight-color:transparent}
.lb-action:hover{color:#e8e6dc}
.lb-img-wrap{display:flex;align-items:center;justify-content:center;flex:1;padding:56px 0 40px;overflow:hidden;position:relative;width:100%}
.lb-img-wrap img{max-width:100%;max-height:100%;object-fit:contain;user-select:none;-webkit-user-drag:none;touch-action:pan-y}
.lb-counter{position:absolute;bottom:12px;left:0;right:0;text-align:center;font-size:13px;color:#73726c;z-index:10;pointer-events:none}
.lb-arrows{position:absolute;top:0;left:0;right:0;bottom:0;display:flex;align-items:center;justify-content:space-between;pointer-events:none;z-index:5;padding:0 4px}
.lb-arrow{pointer-events:auto;background:rgba(0,0,0,.5);border:1px solid rgba(255,255,255,.08);color:#e8e6dc;cursor:pointer;font-size:24px;width:44px;height:44px;border-radius:50%;display:flex;align-items:center;justify-content:center;transition:background .15s,opacity .15s;opacity:.7;-webkit-tap-highlight-color:transparent;touch-action:manipulation;user-select:none}
.lb-arrow:hover{background:rgba(0,0,0,.7);opacity:1}
.lb-arrow:active{opacity:1}
@media(max-width:640px){
  .lb-arrow{width:40px;height:40px;font-size:20px;opacity:.9}
  .lb-img-wrap{padding:48px 0 36px}
}

/* footer */
.footer{margin-top:auto;padding:32px 24px;text-align:center;font-size:12px;color:#5e5d59;line-height:1.6;max-width:1280px;width:100%;margin-left:auto;margin-right:auto}
.footer a{color:#73726c;text-decoration:none}
.footer a:hover{color:#9a9990}

@media(max-width:640px){
  .grid{grid-template-columns:repeat(auto-fill,minmax(140px,1fr));gap:8px;padding:0 12px}
  .header{padding:24px 12px 12px}
  .header h1{font-size:1.35rem}
  .footer{padding:24px 16px;font-size:11px}
}
</style>
</head>
<body>
<div class="header">
  <h1>` + title + `</h1>
  ` + descHTML + `
</div>
<div class="grid" id="gallery">
      ` + gallery + `
</div>

<div class="lightbox" id="lightbox">
  <div class="lightbox-top">
    <span class="lb-name" id="lbName"></span>
    <button class="lb-action" id="lbDl" title="Download"><svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg></button>
    <button class="lb-action" id="lbClose">&times;</button>
  </div>
  <div class="lb-img-wrap" id="lbImgWrap">
    <img id="lbImg" src="" alt="">
    <div class="lb-arrows">
      <button class="lb-arrow" id="lbPrev">&#10094;</button>
      <button class="lb-arrow" id="lbNext">&#10095;</button>
    </div>
  </div>
  <div class="lb-counter" id="lbCounter"></div>
</div>

<div class="footer">
  Phục vụ bởi Caddy và goPanel. Bản quyền goPanel thuộc Lê Hùng Quang Minh.
</div>

<script>
(function(){
  var items=document.querySelectorAll('.thumb-item');
  var images=[],names=[];

  // lazy load images — max 3 concurrent, Intersection Observer triggers load
  var loading=0,maxLoad=3,pending=[];
  function loadNext(){
    while(loading<maxLoad && pending.length){
      var el=pending.shift();
      loadImg(el);
    }
  }
  function loadImg(el){
    var img=el.querySelector('img[data-src]');
    if(!img)return;
    loading++;
    img.src=img.getAttribute('data-src');
    img.removeAttribute('data-src');
    img.onload=function(){img.classList.add('loaded');loading--;loadNext();};
    img.onerror=function(){img.style.display='none';loading--;loadNext();};
  }
  var observer=new IntersectionObserver(function(entries){
    entries.forEach(function(entry){
      if(!entry.isIntersecting)return;
      var el=entry.target;
      var img=el.querySelector('img[data-src]');
      if(!img)return;
      observer.unobserve(el);
      pending.push(el);
      loadNext();
    });
  },{rootMargin:'200px 0px'});

  items.forEach(function(el,i){
    images.push(el.dataset.src);
    names.push(el.querySelector('.thumb-name').textContent);
    el.addEventListener('click',function(){open(i);});
    observer.observe(el);
  });

  var lb=document.getElementById('lightbox');
  var lbImg=document.getElementById('lbImg');
  var lbName=document.getElementById('lbName');
  var lbDl=document.getElementById('lbDl');
  var lbCounter=document.getElementById('lbCounter');
  var lbImgWrap=document.getElementById('lbImgWrap');
  var cur=0;
  var loadedLb={};
  function open(i){
    cur=i;
    if(!loadedLb[i]){
      loadedLb[i]=true;
      lbImg.src=images[i];
    }
    lbName.textContent=names[i];
    lbDl.onclick=function(){var a=document.createElement('a');a.href=images[i];a.download=names[i];a.click();};
    lbCounter.textContent=(i+1)+'/'+images.length;
    lb.classList.add('active');
    document.body.style.overflow='hidden';
  }
  function close(){lb.classList.remove('active');document.body.style.overflow='';}
  function prev(){if(cur>0)open(cur-1);else open(images.length-1);}
  function next(){if(cur<images.length-1)open(cur+1);else open(0);}
  document.getElementById('lbClose').onclick=close;
  document.getElementById('lbPrev').onclick=prev;
  document.getElementById('lbNext').onclick=next;

  // keyboard
  document.addEventListener('keydown',function(e){
    if(!lb.classList.contains('active'))return;
    if(e.key==='Escape')close();
    if(e.key==='ArrowLeft')prev();
    if(e.key==='ArrowRight')next();
  });

  // swipe
  var startX=0,startY=0,swiping=false;
  lbImgWrap.addEventListener('touchstart',function(e){
    var t=e.touches[0];
    startX=t.clientX;startY=t.clientY;swiping=true;
  },{passive:true});
  lbImgWrap.addEventListener('touchmove',function(e){
    if(!swiping)return;
    var t=e.touches[0];
    var dx=t.clientX-startX,dy=t.clientY-startY;
    if(Math.abs(dx)>30 && Math.abs(dx)>Math.abs(dy)*1.5){
      e.preventDefault();
      swiping=false;
      if(dx>0)prev();else next();
    }
  },{passive:false});

  // close on bg click
  lb.addEventListener('click',function(e){
    if(e.target===lb||e.target===lbImgWrap||e.target===lbImg)close();
  });
})();
</script>
</body>
</html>`
}
