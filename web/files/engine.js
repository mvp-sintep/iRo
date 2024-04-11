/*<![CDATA[*/
const A=(arg)=>{if(!arg)return([]);if(arg.toArray)return(arg.toArray());let x=arg.length,a=new Array(x);while(x--)a[x]=arg[x];return(a);};
Object["extend"]=function(arg){
  if(!!arg) {
    for(let key,src,i=1,max=arguments.length;i<max;i++){
      src=arguments[i];
      for(key in src)
        if(src.hasOwnProperty(key)){
          if(!src[key].constructor.extend)arg[key]=src[key];
          else Object.extend(arg[key]={},src[key]);
        }
    }
  }
  return arg;
};
const HTMLElementInterface = {
  extend: function () { return (Object.extend.apply(this,[this].concat(A(arguments)))); },
  destroy: function () { this.remove(); },
  isExist: function () { return (!!this.document); },
  insertTo: function (arg) { arg.insertBefore(this); return this; },
  appendTo: function (arg) { arg.appendChild(this); return this; },
  withID: function (arg) { this.id = arg; return this; },
  withInnerText: function (arg) { this.innerText = arg; return this; },
  withInnerHTML: function (arg) { this.innerHTML = arg; return this; },
  withClassName: function (arg) { this.className = arg; return this; },
  withAttribute: function (key,value) { this.setAttribute(key,value); return this; }
};
const one = {tick:55,second:1000,minute:60000,hour:3600000,date:86400000};
const DateExpandingInterface = {
  truncYear: function (arg = 1) { return new Date(this.getFullYear() - this.getFullYear() % arg, 0, 1); },
  withYearOffset: function (arg = 1) { return (new Date(this.getFullYear() + arg, this.getMonth(), this.getDate(), this.getHours(), this.getMinutes(), this.getSeconds(), this.getMilliseconds())); },
  getYearDuration: function () { let tmp = this.truncYear(); return (tmp.withYearOffset(1).valueOf() - tmp.valueOf()); },
  truncMonth: function (arg = 1) { return (new Date(this.getFullYear(), this.getMonth() - this.getMonth() % arg, 1)); },
  withMonthOffset: function (arg = 1) { return (new Date(this.getFullYear(), this.getMonth() + arg, this.getDate(), this.getHours(), this.getMinutes(), this.getSeconds(), this.getMilliseconds())); },
  getMonthDuration: function () { let tmp = this.truncMonth(); return (tmp.withMonthOffset(1).valueOf() - tmp.valueOf()); },
  truncDate: function (arg = 1) { return (new Date(this.getFullYear(), this.getMonth(), this.getDate() - this.getDate() % arg)); },
  truncHours: function (arg = 1) { return (new Date(this.getFullYear(), this.getMonth(), this.getDate(), this.getHours() - this.getHours() % arg)); },
  truncMinutes: function (arg = 1) { return (new Date(this.getFullYear(), this.getMonth(), this.getDate(), this.getHours(), this.getMinutes() - this.getMinutes() % arg)); },
  truncSeconds: function (arg = 1) { return (new Date(this.getFullYear(), this.getMonth(), this.getDate(), this.getHours(), this.getMinutes(), this.getSeconds() - this.getSeconds() % arg)); },
  withOffset: function (arg = 0) { return (new Date(this.valueOf() + arg)); },
  getOffset: function (arg = new Date()) { let result = 0 - this.valueOf(); try { result += arg.valueOf(); } catch (e) { } return new Date(result); }
};
Object.extend(Date.prototype, DateExpandingInterface);
/*]]>*/
