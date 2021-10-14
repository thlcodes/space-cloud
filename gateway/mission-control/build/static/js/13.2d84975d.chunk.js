(this["webpackJsonpmission-control"]=this["webpackJsonpmission-control"]||[]).push([[13],{1458:function(e,a,t){"use strict";t.r(a);t(575);var n=t(574),r=(t(581),t(582)),l=(t(572),t(573)),c=(t(593),t(591)),i=(t(588),t(586)),s=(t(632),t(646)),u=(t(562),t(563)),m=(t(604),t(610)),o=(t(548),t(547)),d=(t(328),t(209)),g=(t(580),t(578)),p=(t(549),t(550)),y=(t(571),t(569)),E=t(709),h=t.n(E),v=t(836),f=(t(551),t(552)),k=t(53),b=t(0),P=t.n(b),S=t(1474),w=t(1472),I=t(1487),C=t(553),K=t(13),A=t(567),R=t(1433),U=t.n(R),j=t(1434),L=t.n(j),T=t(28),O=function(e){var a=e.handleSubmit,t=e.handleCancel,n=e.initialValues,r=f.a.useForm(),l=Object(k.a)(r,1)[0],c=Object(b.useState)(!1),i=Object(k.a)(c,2),s=i[0],E=i[1],I={alg:n?n.alg:"HS256",secret:n?n.secret:Object(K.j)(),privateKey:n?n.privateKey:void 0,publicKey:n?n.publicKey:void 0,jwkUrl:n?n.jwkUrl:void 0,isPrimary:!!n&&n.isPrimary,checkAudience:!!(n&&n.aud&&n.aud.length>0),checkIssuer:!(!n||!n.iss),aud:n&&n.aud?n.aud:[""],iss:n&&n.iss?n.iss:[""],kid:n&&n.kid?n.kid:Object(T.generateId)()};return P.a.createElement(u.a,{title:"Add secret",okText:n?"Save":"Add",visible:!0,onOk:function(e){l.validateFields().then((function(e){var n=e=Object.assign({},I,e),r=n.alg,l=n.secret,c=n.privateKey,i=n.publicKey,s=n.jwkUrl,u=n.isPrimary,m=n.checkAudience,o=n.checkIssuer,d=n.aud,g=n.iss,p=n.kid;"RS256"===r&&c&&(i=function(e){var a=L.a.pki.privateKeyFromPem(e),t=L.a.pki.setRsaPublicKey(a.n,a.e);return L.a.pki.publicKeyToPem(t)}(c)),m||(d=void 0),o||(g=void 0);var y={};switch(r){case"HS256":y={alg:r,secret:l,isPrimary:u,aud:d,iss:g,kid:p};break;case"RS256":y={alg:r,publicKey:i,privateKey:c,isPrimary:u,aud:d,iss:g,kid:p};break;case"RS256_PUBLIC":y={alg:r,publicKey:i,aud:d,iss:g,kid:p};break;case"JWK_URL":y={alg:r,jwkUrl:s,aud:d,iss:g}}a(y).then((function(){return t()}))}))},onCancel:t,width:720},P.a.createElement(f.a,{form:l,layout:"vertical",initialValues:I},P.a.createElement(C.a,{name:"Algorithm"}),P.a.createElement(f.a.Item,{name:"alg"},P.a.createElement(y.a,{onChange:function(e){"RS256"===e&&(E(!0),setTimeout(Object(v.a)(h.a.mark((function e(){var a;return h.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.next=2,U()().private;case 2:a=e.sent,l.setFieldsValue({privateKey:a}),E(!1);case 5:case"end":return e.stop()}}),e)}))),500))}},P.a.createElement(y.a.Option,{value:"HS256"},"HS256"),P.a.createElement(y.a.Option,{value:"RS256"},"RS256"),P.a.createElement(y.a.Option,{value:"RS256_PUBLIC"},"RS256 PUBLIC"),P.a.createElement(y.a.Option,{value:"JWK_URL"},"JWK URL"))),P.a.createElement(A.a,{dependency:"alg",condition:function(){return"HS256"===l.getFieldValue("alg")}},P.a.createElement(C.a,{name:"Secret"}),P.a.createElement(f.a.Item,{name:"secret",rules:[{required:!0,message:"Please provide a secret"}]},P.a.createElement(p.a.Password,{className:"input",placeholder:"Secret value"})),P.a.createElement(C.a,{name:"Primary secret"}),P.a.createElement(f.a.Item,{name:"isPrimary",valuePropName:"checked"},P.a.createElement(g.a,null,"Use this secret in user management module of API gateway to sign tokens on successful signup/signin requests"))),P.a.createElement(A.a,{dependency:"alg",condition:function(){return"RS256"===l.getFieldValue("alg")}},P.a.createElement(C.a,{name:"Private key"}),P.a.createElement(f.a.Item,{name:"privateKey"},s?P.a.createElement("span",null,P.a.createElement(d.a,{className:"page-loading",spinning:!0,size:"large"})," Generating Private Key......"):P.a.createElement(p.a.TextArea,{rows:4,placeholder:"Private key"})),P.a.createElement(C.a,{name:"Primary secret"}),P.a.createElement(f.a.Item,{name:"isPrimary",valuePropName:"checked"},P.a.createElement(g.a,null,"Use this secret in user management module of API gateway to sign tokens on successful signup/signin requests"))),P.a.createElement(A.a,{dependency:"alg",condition:function(){return"JWK_URL"===l.getFieldValue("alg")}},P.a.createElement(C.a,{name:"JWK URL"}),P.a.createElement(f.a.Item,{name:"jwkUrl",rules:[{required:!0,message:"Please provide an URL"}]},P.a.createElement(p.a,{placeholder:"JWK URL"}))),P.a.createElement(A.a,{dependency:"alg",condition:function(){return"RS256_PUBLIC"===l.getFieldValue("alg")}},P.a.createElement(C.a,{name:"Public key"}),P.a.createElement(f.a.Item,{name:"publicKey"},P.a.createElement(p.a.TextArea,{rows:4,placeholder:"Public key"}))),P.a.createElement(m.a,{className:"advanced",style:{background:"white"},bordered:!1},P.a.createElement(m.a.Panel,{header:"Advanced",key:"advanced"},P.a.createElement(C.a,{name:"Check audience"}),P.a.createElement(f.a.Item,{name:"checkAudience",valuePropName:"checked"},P.a.createElement(g.a,null,"Check audience while verifying JWT token")),P.a.createElement(A.a,{dependency:"checkAudience",condition:function(){return!0===l.getFieldValue("checkAudience")}},P.a.createElement(C.a,{name:"Specify audiences",description:"The audience check will pass if the JWT matches any one of the specified audiences below"}),P.a.createElement(f.a.List,{name:"aud"},(function(e,a){var t=a.add,n=a.remove;return P.a.createElement("div",null,e.map((function(a){return P.a.createElement(f.a.Item,{key:a.key,style:{marginBottom:8}},P.a.createElement(f.a.Item,Object.assign({},a,{validateTrigger:["onChange","onBlur"],rules:[{required:!0,message:"Please input a value"}],noStyle:!0}),P.a.createElement(p.a,{placeholder:"Audience value",style:{width:"90%"}})),e.length>1?P.a.createElement(S.a,{style:{marginLeft:16},onClick:function(){n(a.name)}}):null)})),P.a.createElement(f.a.Item,null,P.a.createElement(o.a,{onClick:function(){l.validateFields(e.map((function(e){return["aud",e.name]}))).then((function(){return t()})).catch((function(e){return console.log("Exception",e)}))}},P.a.createElement(w.a,null)," Add")))}))),P.a.createElement(C.a,{name:"Check issuer"}),P.a.createElement(f.a.Item,{name:"checkIssuer",valuePropName:"checked"},P.a.createElement(g.a,null,"Check issuer while verifying JWT token")),P.a.createElement(A.a,{dependency:"checkIssuer",condition:function(){return!0===l.getFieldValue("checkIssuer")}},P.a.createElement(C.a,{name:"Specify issuers",description:"The issuer check will pass if the JWT matches any one of the specified issuers below"}),P.a.createElement(f.a.List,{name:"iss"},(function(e,a){var t=a.add,n=a.remove;return P.a.createElement("div",null,e.map((function(a){return P.a.createElement(f.a.Item,{key:a.key,style:{marginBottom:8}},P.a.createElement(f.a.Item,Object.assign({},a,{validateTrigger:["onChange","onBlur"],rules:[{required:!0,message:"Please input a value"}],noStyle:!0}),P.a.createElement(p.a,{placeholder:"Issuer value",style:{width:"90%"}})),e.length>1?P.a.createElement(S.a,{style:{marginLeft:16},onClick:function(){n(a.name)}}):null)})),P.a.createElement(f.a.Item,null,P.a.createElement(o.a,{onClick:function(){l.validateFields(e.map((function(e){return["iss",e.name]}))).then((function(){return t()})).catch((function(e){return console.log("Exception",e)}))}},P.a.createElement(w.a,null)," Add")))}))),P.a.createElement(A.a,{dependency:"alg",condition:function(){return"JWK_URL"!==l.getFieldValue("alg")}},P.a.createElement(C.a,{name:"kid",hint:"(key ID)"}),P.a.createElement(f.a.Item,{name:"kid",rules:[{required:!0,message:"Please input a value!"}]},P.a.createElement(p.a,{placeholder:"kid"})))))))},J=function(e){var a=e.secretData,t=e.handleCancel;return P.a.createElement(u.a,{title:"View secret",visible:!0,footer:null,onCancel:t,width:720},P.a.createElement(s.a.Paragraph,{style:{fontSize:16},strong:!0},"Algorithm"),P.a.createElement("div",{style:{marginBottom:24}},a.alg),"HS256"===a.alg&&P.a.createElement(P.a.Fragment,null,P.a.createElement(s.a.Paragraph,{style:{fontSize:16},copyable:{text:a.secret},strong:!0},"Secret"),P.a.createElement(p.a,{value:a.secret})),"RS256"===a.alg&&P.a.createElement(P.a.Fragment,null,P.a.createElement(s.a.Paragraph,{style:{fontSize:16},copyable:{text:a.publicKey},strong:!0},"Public key"),P.a.createElement("div",{style:{marginBottom:24}},P.a.createElement(p.a.TextArea,{rows:4,value:a.publicKey})),P.a.createElement(s.a.Paragraph,{style:{fontSize:16},copyable:{text:a.privateKey},strong:!0},"Private key"),P.a.createElement("div",{style:{marginBottom:24}},P.a.createElement(p.a.TextArea,{rows:4,value:a.privateKey}))),"JWK_URL"===a.alg&&P.a.createElement(P.a.Fragment,null,P.a.createElement(s.a.Paragraph,{style:{fontSize:16},copyable:{text:a.jwkUrl},strong:!0},"URL"),P.a.createElement(p.a,{value:a.jwkUrl})),"RS256_PUBLIC"===a.alg&&P.a.createElement(P.a.Fragment,null,P.a.createElement(s.a.Paragraph,{style:{fontSize:16},copyable:{text:a.publicKey},strong:!0},"Public key"),P.a.createElement("div",{style:{marginBottom:24}},P.a.createElement(p.a.TextArea,{rows:4,value:a.publicKey}))))};a.default=function(e){var a=e.secrets,t=e.handleRemoveSecret,s=e.handleChangePrimarySecret,u=e.handleSaveSecret,m=a.map((function(e){return Object.assign({},e,{alg:e.alg?e.alg:"HS256"})})),d=Object(b.useState)(!1),g=Object(k.a)(d,2),p=g[0],y=g[1],E=Object(b.useState)(!1),h=Object(k.a)(E,2),v=h[0],f=h[1],S=Object(b.useState)(void 0),w=Object(k.a)(S,2),C=w[0],K=w[1],A=void 0!==C?m[C]:null,R=[{title:"Algorithm",dataIndex:"alg"},{title:P.a.createElement("span",null,"Primary secret  ",P.a.createElement(i.a,{placement:"bottomLeft",title:"Primary secret is used by the user management module of API gateway to sign tokens on successful signup/signin requests."},P.a.createElement(I.a,null))),render:function(e,a,t){return"JWK_URL"===a.alg||"RS256_PUBLIC"===a.alg?P.a.createElement("span",null,"N/A"):P.a.createElement(c.a,{checked:a.isPrimary,onChange:a.isPrimary?void 0:function(){return s(t)}})}},{title:"Actions",className:"column-actions",render:function(e,a,n){return P.a.createElement("span",null,P.a.createElement("a",{onClick:function(){return function(e){K(e),f(!0)}(n)}},"View"),P.a.createElement("a",{onClick:function(){return function(e){K(e),y(!0)}(n)}},"Edit"),P.a.createElement(l.a,{title:a.isPrimary?"You are deleting primary secret. Any remaining secret will be randomly chosen as primary key.":"Tokens signed with this secret will stop getting verified",onConfirm:function(){return t(n)}},P.a.createElement("a",{style:{color:"red"}},"Remove")))}}];return P.a.createElement("div",null,P.a.createElement("div",{style:{display:"flex",justifyContent:"space-between"}},P.a.createElement("h2",{style:{display:"inline-block"}},"JWT Secrets"),P.a.createElement(o.a,{onClick:function(){return y(!0)}},"Add")),P.a.createElement("p",null,"These secrets are used by the auth module in Space Cloud to verify the JWT token for all API requests"),P.a.createElement(r.a,{description:"Space Cloud supports multiple JWT secrets to safely rotate them",type:"info",showIcon:!0}),P.a.createElement(n.a,{style:{marginTop:16},columns:R,dataSource:m,bordered:!0,pagination:!1,rowKey:"key"}),p&&P.a.createElement(O,{initialValues:A,handleSubmit:function(e){return u(e,C)},handleCancel:function(){y(!1),K(void 0)}}),v&&P.a.createElement(J,{secretData:A,handleCancel:function(){f(!1),K(void 0)}}))}},794:function(e,a){}}]);
//# sourceMappingURL=13.2d84975d.chunk.js.map