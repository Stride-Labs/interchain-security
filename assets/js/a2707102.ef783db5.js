"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[543],{3905:(e,t,r)=>{r.d(t,{Zo:()=>d,kt:()=>h});var n=r(7294);function a(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function i(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function o(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?i(Object(r),!0).forEach((function(t){a(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):i(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function l(e,t){if(null==e)return{};var r,n,a=function(e,t){if(null==e)return{};var r,n,a={},i=Object.keys(e);for(n=0;n<i.length;n++)r=i[n],t.indexOf(r)>=0||(a[r]=e[r]);return a}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(n=0;n<i.length;n++)r=i[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(a[r]=e[r])}return a}var c=n.createContext({}),p=function(e){var t=n.useContext(c),r=t;return e&&(r="function"==typeof e?e(t):o(o({},t),e)),r},d=function(e){var t=p(e.components);return n.createElement(c.Provider,{value:t},e.children)},s="mdxType",u={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},m=n.forwardRef((function(e,t){var r=e.components,a=e.mdxType,i=e.originalType,c=e.parentName,d=l(e,["components","mdxType","originalType","parentName"]),s=p(r),m=a,h=s["".concat(c,".").concat(m)]||s[m]||u[m]||i;return r?n.createElement(h,o(o({ref:t},d),{},{components:r})):n.createElement(h,o({ref:t},d))}));function h(e,t){var r=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var i=r.length,o=new Array(i);o[0]=m;var l={};for(var c in t)hasOwnProperty.call(t,c)&&(l[c]=t[c]);l.originalType=e,l[s]="string"==typeof e?e:a,o[1]=l;for(var p=2;p<i;p++)o[p]=r[p];return n.createElement.apply(null,o)}return n.createElement.apply(null,r)}m.displayName="MDXCreateElement"},1663:(e,t,r)=>{r.r(t),r.d(t,{assets:()=>c,contentTitle:()=>o,default:()=>u,frontMatter:()=>i,metadata:()=>l,toc:()=>p});var n=r(7462),a=(r(7294),r(3905));const i={sidebar_position:1,title:"ADRs"},o="Architecture Decision Records (ADR)",l={unversionedId:"adrs/intro",id:"adrs/intro",title:"ADRs",description:"This is a location to record all high-level architecture decisions in the Interchain Security project.",source:"@site/docs/adrs/intro.md",sourceDirName:"adrs",slug:"/adrs/intro",permalink:"/interchain-security/adrs/intro",draft:!1,tags:[],version:"current",sidebarPosition:1,frontMatter:{sidebar_position:1,title:"ADRs"},sidebar:"tutorialSidebar",previous:{title:"Frequently Asked Questions",permalink:"/interchain-security/faq"},next:{title:"ADR Template",permalink:"/interchain-security/adrs/adr-007-pause-unbonding-on-eqv-prop"}},c={},p=[{value:"Table of Contents",id:"table-of-contents",level:2}],d={toc:p},s="wrapper";function u(e){let{components:t,...r}=e;return(0,a.kt)(s,(0,n.Z)({},d,r,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("h1",{id:"architecture-decision-records-adr"},"Architecture Decision Records (ADR)"),(0,a.kt)("p",null,"This is a location to record all high-level architecture decisions in the Interchain Security project."),(0,a.kt)("p",null,"You can read more about the ADR concept in this ",(0,a.kt)("a",{parentName:"p",href:"https://product.reverb.com/documenting-architecture-decisions-the-reverb-way-a3563bb24bd0#.78xhdix6t"},"blog post"),"."),(0,a.kt)("p",null,"An ADR should provide:"),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},"Context on the relevant goals and the current state"),(0,a.kt)("li",{parentName:"ul"},"Proposed changes to achieve the goals"),(0,a.kt)("li",{parentName:"ul"},"Summary of pros and cons"),(0,a.kt)("li",{parentName:"ul"},"References"),(0,a.kt)("li",{parentName:"ul"},"Changelog")),(0,a.kt)("p",null,"Note the distinction between an ADR and a spec. The ADR provides the context, intuition, reasoning, and\njustification for a change in architecture, or for the architecture of something\nnew. The spec is much more compressed and streamlined summary of everything as\nit is or should be."),(0,a.kt)("p",null,"If recorded decisions turned out to be lacking, convene a discussion, record the new decisions here, and then modify the code to match."),(0,a.kt)("p",null,"Note the context/background should be written in the present tense."),(0,a.kt)("p",null,"To suggest an ADR, please make use of the ",(0,a.kt)("a",{parentName:"p",href:"/interchain-security/adrs/adr-template"},"ADR template")," provided."),(0,a.kt)("h2",{id:"table-of-contents"},"Table of Contents"),(0,a.kt)("table",null,(0,a.kt)("thead",{parentName:"table"},(0,a.kt)("tr",{parentName:"thead"},(0,a.kt)("th",{parentName:"tr",align:null},"ADR ","#"),(0,a.kt)("th",{parentName:"tr",align:null},"Description"),(0,a.kt)("th",{parentName:"tr",align:null},"Status"))),(0,a.kt)("tbody",{parentName:"table"},(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("a",{parentName:"td",href:"/interchain-security/adrs/adr-001-key-assignment"},"001")),(0,a.kt)("td",{parentName:"tr",align:null},"Consumer chain key assignment"),(0,a.kt)("td",{parentName:"tr",align:null},"Accepted, Implemented")),(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("a",{parentName:"td",href:"/interchain-security/adrs/adr-002-throttle"},"002")),(0,a.kt)("td",{parentName:"tr",align:null},"Jail Throttling"),(0,a.kt)("td",{parentName:"tr",align:null},"Accepted, Implemented")),(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("a",{parentName:"td",href:"/interchain-security/adrs/adr-003-equivocation-gov-proposal"},"003")),(0,a.kt)("td",{parentName:"tr",align:null},"Equivocation governance proposal"),(0,a.kt)("td",{parentName:"tr",align:null},"Accepted, Implemented")),(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("a",{parentName:"td",href:"./adr-004-denom-dos-fixes"},"004")),(0,a.kt)("td",{parentName:"tr",align:null},"Denom DOS fixes"),(0,a.kt)("td",{parentName:"tr",align:null},"Accepted, Implemented")),(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("a",{parentName:"td",href:"/interchain-security/adrs/adr-005-cryptographic-equivocation-verification"},"005")),(0,a.kt)("td",{parentName:"tr",align:null},"Cryptographic verification of equivocation evidence"),(0,a.kt)("td",{parentName:"tr",align:null},"Accepted, In-progress")),(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("a",{parentName:"td",href:"/interchain-security/adrs/adr-007-pause-unbonding-on-eqv-prop"},"007")),(0,a.kt)("td",{parentName:"tr",align:null},"Pause validator unbonding during equivocation proposal"),(0,a.kt)("td",{parentName:"tr",align:null},"Proposed")),(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("a",{parentName:"td",href:"/interchain-security/adrs/adr-008-throttle-retries"},"008")),(0,a.kt)("td",{parentName:"tr",align:null},"Throttle with retries"),(0,a.kt)("td",{parentName:"tr",align:null},"Accepted, In-progress")),(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("a",{parentName:"td",href:"/interchain-security/adrs/adr-009-soft-opt-out"},"009")),(0,a.kt)("td",{parentName:"tr",align:null},"Soft Opt-out"),(0,a.kt)("td",{parentName:"tr",align:null},"Accepted, Implemented")),(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("a",{parentName:"td",href:"/interchain-security/adrs/adr-010-standalone-changeover"},"009")),(0,a.kt)("td",{parentName:"tr",align:null},"Standalone to Consumer Changeover"),(0,a.kt)("td",{parentName:"tr",align:null},"Accepted, Implemented")))))}u.isMDXComponent=!0}}]);