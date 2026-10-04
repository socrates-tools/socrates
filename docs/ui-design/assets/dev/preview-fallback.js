// Preview helper: resolves the compiled design-system namespace, or - before the
// bundle has been compiled - evaluates the component sources directly with Babel.
(function () {
  var here = document.currentScript && document.currentScript.src;
  var root = here ? new URL('../../', here).href : '../../';
  var ORDER = ['brand/Icon', 'marks/Tag', 'marks/Check', 'marks/Mark', 'brand/Wordmark', 'marks/StatusPill', 'actions/Button',
    'forms/Option', 'forms/TextArea', 'forms/Progress', 'evidence/CodeBlock', 'evidence/DiffView', 'evidence/Illustration',
    'surfaces/Card', 'surfaces/Notice', 'navigation/IndexEntry', 'navigation/TopBar', 'feedback/Dialog', 'feedback/Tooltip'];
  function compiled() {
    var k = Object.keys(window).find(function (n) { try { var v = window[n]; return v && typeof v === 'object' && v.Mark && v.Option && v.DiffView; } catch (e) { return false; } });
    return k ? window[k] : null;
  }
  window.__soReady = (async function () {
    var ns = compiled();
    if (ns) return ns;
    var so = {};
    for (var i = 0; i < ORDER.length; i++) {
      var src = await (await fetch(root + 'components/' + ORDER[i] + '.jsx')).text();
      var names = []; src.replace(/export function (\w+)/g, function (_, n) { names.push(n); });
      src = src.replace(/^import .*$/gm, '').replace(/export function/g, 'function');
      var have = Object.keys(so);
      var wrapped = '(function(so){' + (have.length ? 'var {' + have.join(',') + '} = so;' : '') + src + '; return {' + names.join(',') + '};})';
      var fn = (0, eval)(Babel.transform(wrapped, { presets: ['react'] }).code);
      Object.assign(so, fn(so));
    }
    return so;
  })();
})();
