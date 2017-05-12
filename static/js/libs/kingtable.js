/**
 * KingTable 2.0.0
 * https://github.com/RobertoPrevato/KingTable
 *
 * Copyright 2017, Roberto Prevato
 * https://robertoprevato.github.io
 *
 * Licensed under the MIT license:
 * http://www.opensource.org/licenses/MIT
 */
! function e(t, n, r) {
	function i(s, o) {
		if (!n[s]) {
			if (!t[s]) {
				var u = "function" == typeof require && require;
				if (!o && u) return u(s, !0);
				if (a) return a(s, !0);
				var l = new Error("Cannot find module '" + s + "'");
				throw l.code = "MODULE_NOT_FOUND", l
			}
			var f = n[s] = {
				exports: {}
			};
			t[s][0].call(f.exports, function (e) {
				var n = t[s][1][e];
				return i(n || e)
			}, f, f.exports, e, t, n, r)
		}
		return n[s].exports
	}
	for (var a = "function" == typeof require && require, s = 0; s < r.length; s++) i(r[s]);
	return i
}({
	1: [function (e, t, n) {
		"use strict";

		function r(e) {
			return e && e.__esModule ? e : {
				default: e
			}
		}

		function i(e) {
			for (var t = e.length, n = 0; n < t; n++) {
				var r = e[n],
					i = r[0],
					a = r[1];
				u.default.isNumber(a) || /^asc|^desc/i.test(a) || (0, m.ArgumentException)("The sort order '" + a + "' for '" + i + "' is not valid (it must be /^asc|^desc/i)."), r[1] = u.default.isNumber(a) ? a : /^asc/i.test(a) ? 1 : -1
			}
			return e
		}

		function a(e) {
			var t = e.split(/([\.\,]\d+)$/),
				n = t[0],
				r = t[1],
				i = 0;
			return n && (i = parseInt(n.replace(/\D/g, ""))), r && (i += parseFloat(r.replace(/\,/g, "."))), /^\s?-/.test(e) ? -i : i
		}

		function s(e) {
			u.default.isString(e) || (0, m.TypeException)("s", "string");
			var t = e.match(/[-+~]?([0-9]{1,3}(?:[,\s\.]{1}[0-9]{3})*(?:[\.|\,]{1}[0-9]+)?)/g);
			if (t && 1 == t.length) {
				if (/(#[0-9a-fA-F]{3}|#[0-9a-fA-F]{6}|#[0-9a-fA-F]{8})$/.test(e)) return !1;
				//if (e.match(/[^0-9\.\,\s]/g).length > 6) return !1;
				return a(t[0])
			}
			return !1
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var o = e("../../scripts/utils"),
			u = r(o),
			l = e("../../scripts/components/string"),
			f = r(l),
			c = e("../../scripts/components/regex"),
			d = r(c),
			p = e("../../scripts/components/reflection"),
			h = r(p),
			m = e("../../scripts/exceptions"),
			g = {
				autoParseNumbers: !0,
				ci: !0
			};
		n.default = {
			normalizeOrder: i,
			lookSortableAsNumber: s,
			options: g,
			parseSortBy: function (e) {
				if (e) {
					var t = e.split(/\s*,\s*/g);
					return u.default.map(t, function (e) {
						var t = e.split(/\s/),
							n = t[0],
							r = t[1] || "asc";
						return [n, f.default.startsWith(r, "asc", !0) ? 1 : -1]
					})
				}
			},
			humanSortBy: function (e, t) {
				return e && e.length ? u.default.map(e, function (e) {
					var n = e[0];
					return 1 === e[1] ? t ? n + " asc" : n : n + " desc"
				}).join(", ") : ""
			},
			getSortCriteria: function (e) {
				var t, n = this,
					r = e.length;
				if (1 == e.length) {
					var i = e[0];
					if (u.default.isString(i) && i.search(/,|\s/) > -1) return this.parseSortBy(i)
				}
				if (r > 1) {
					var a = u.default.toArray(e);
					t = u.default.map(a, function (e) {
						return n.normalizeSortByValue(e, !0)
					})
				} else t = this.normalizeSortByValue(e[0]);
				return t
			},
			normalizeSortByValue: function (e, t) {
				var n = this;
				if (u.default.isString(e)) return t ? [e, "asc"] : [[e, "asc"]];
				if (u.default.isArray(e)) return u.default.isArray(e[0]) ? e : u.default.map(e, function (e) {
					return n.normalizeSortByValue(e, !0)
				});
				if (u.default.isPlainObject(e)) {
					var r, i = [];
					for (r in e) i.push([r, e[r]]);
					return i
				}(0, m.TypeException)("sort", "string | [] | {}")
			},
			compareStrings: function (e, t, n) {
				if (this.options.autoParseNumbers) {
					var r = s(e),
						i = s(t);
					if (!1 !== r || !1 !== i) {
						if (r === i) return 0;
						if (!1 !== r && !1 === t) return n;
						if (!1 === r && !1 !== t) return -n;
						if (r < i) return -n;
						if (r > i) return n
					}
				}
				return f.default.compare(e, t, n, this.options)
			},
			sortBy: function (e) {
				u.default.isArray(e) || (0, m.TypeException)("ar", "array");
				var t = arguments.length,
					n = u.default.toArray(arguments).slice(1, t),
					r = this.getSortCriteria(n);
				r = i(r);
				var a = r.length,
					s = u.default.isString,
					o = this.compareStrings.bind(this),
					l = void 0,
					f = null;
				return e.sort(function (e, t) {
					if (e === t) return 0;
					if (e !== l && t === l) return -1;
					if (e === l && t !== l) return 1;
					if (e !== f && t === f) return -1;
					if (e === f && t !== f) return 1;
					for (var n = 0; n < a; n++) {
						var i = r[n],
							u = i[0],
							c = i[1],
							d = e[u],
							p = t[u];
						if (d !== p) {
							if (d !== l && p === l) return -c;
							if (d === l && p !== l) return c;
							if (d !== f && p === f) return -c;
							if (d === f && p !== f) return c;
							if (d && !p) return c;
							if (!d && p) return -c;
							if (s(d) && s(p)) return o(d, p, c);
							if (d < p) return -c;
							if (d > p) return c
						}
					}
					return 0
				}), e
			},
			sortByProperty: function (e, t, n) {
				u.default.isArray(e) || (0, m.TypeException)("arr", "array"), u.default.isString(t) || (0, m.TypeException)("property", "string"), u.default.isUnd(n) || (n = "asc"), n = u.default.isNumber(n) ? n : /^asc/i.test(n) ? 1 : -1;
				var r = {};
				return r[t] = n, this.sortBy(e, r)
			},
			searchByStringProperty: function (e) {
				return u.default.require(e, ["pattern", "collection", "property"]), this.searchByStringProperties(u.default.extend(e, {
					properties: [e.property]
				}))
			},
			search: function (e) {
				if (!e || !e.length) return e;
				var t = arguments.length;
				if (t < 2) return e;
				for (var n, r, i = u.default.toArray(arguments).slice(1, t), a = [], s = u.default.isString, o = 0, t = e.length; o < t; o++) {
					r = e[o];
					for (n in r) s(r[n]) && -1 == a.indexOf(n) && a.push(n)
				}
				return u.default.each(i, function (e) {
					u.default.isString(e) || (0, m.ArgumentException)("Unexpected parameter " + e)
				}), this.searchByStringProperties({
					collection: e,
					pattern: d.default.getPatternFromStrings(i),
					properties: a,
					normalize: !0
				})
			},
			searchByStringProperties: function (e) {
				var t = {
						order: "asc",
						limit: null,
						keepSearchDetails: !1,
						getResults: function (e) {
							if (this.decorate)
								for (var t = 0, n = e.length; t < n; t++) {
									var r = e[t],
										i = u.default.where(r.matches, function (e) {
											return null != e
										});
									e[t].obj.__search_matches__ = i.length ? u.default.map(i, function (e) {
										return e.matchedProperty
									}) : []
								}
							if (this.keepSearchDetails) return e;
							for (var a = [], t = 0, n = e.length; t < n; t++) a.push(e[t].obj);
							return a
						},
						normalize: !0,
						decorate: !1
					},
					n = u.default.extend({}, t, e);
				n.order && n.order.match(/asc|ascending|desc|descending/i) || (n.order = "asc");
				var r = [],
					i = n.pattern;
				if (!(i instanceof RegExp)) {
					if (!u.default.isString(i)) throw new Error("the pattern must be a string or a regular expression");
					i = d.default.getSearchPattern(i)
				}
				for (var a = n.properties, s = "length", o = n.normalize, l = n.collection, c = u.default.isArray, p = u.default.isNumber, m = u.default.flatten, g = (u.default.map, 0), v = l[s]; g < v; g++) {
					for (var y = l[g], b = [], w = 0, k = 0, x = a[s]; k < x; k++) {
						var E = a[k],
							S = h.default.getPropertyValue(y, E);
						if (S)
							if (S.match || (S = S.toString()), c(S)) {
								if (!S[s]) continue;
								S = m(S);
								for (var P, T = [], _ = 0, v = S[s]; _ < v; _++) {
									var C = S[_].match(i);
									C && (p(P) || (P = _), T.push(C))
								}
								T[s] && (b[k] = {
									matchedProperty: E,
									indexes: [P],
									recourrences: m(T)[s]
								})
							} else {
								o && (S = f.default.normalize(S));
								var C = S.match(i);
								if (C) {
									for (var O, V = new RegExp(i.source, "gi"), j = []; O = V.exec(S);) j.push(O.index);
									w += C[s], b[k] = {
										matchValue: S,
										matchedProperty: E,
										indexes: j,
										recourrences: C[s]
									}
								}
							}
					}
					b[s] && r.push({
						obj: y,
						matches: b,
						totalMatches: w
					})
				}
				var M = n.order.match(/asc|ascending/i) ? 1 : -1,
					F = "toLowerCase",
					D = "toString",
					H = "matchedProperty",
					I = "indexOf",
					N = "hasOwnProperty",
					A = "recourrences",
					y = "obj",
					L = "indexes",
					R = "totalMatches";
				r.sort(function (e, t) {
					if (e[R] > t[R]) return -M;
					if (e[R] < t[R]) return M;
					for (var n = 0, r = a[s]; n < r; n++) {
						var i = e.matches[n],
							o = t.matches[n];
						if (i || o) {
							if (i && !o) return -M;
							if (!i && o) return M;
							var l = u.default.min(i[L]),
								f = u.default.min(o[L]);
							if (l < f) return -M;
							if (l > f) return M;
							if (i[L][I](l) < o[L][I](f)) return -M;
							if (i[L][I](l) > o[L][I](f)) return M;
							var c = e[y],
								d = t[y];
							if (c[N](i[H]) && d[N](o[H])) {
								if (c[i[H]][D]()[F]() < d[o[H]][D]()[F]()) return -M;
								if (c[i[H]][D]()[F]() > d[o[H]][D]()[F]()) return M
							}
							if (i[A] > o[A]) return -M;
							if (i[A] < o[A]) return M
						}
					}
					return 0
				});
				var B = n.limit;
				return B && (r = r.slice(0, u.default.min(B, r[s]))), n.getResults(r)
			}
		}
	}, {
		"../../scripts/components/reflection": 5,
		"../../scripts/components/regex": 6,
		"../../scripts/components/string": 7,
		"../../scripts/exceptions": 20,
		"../../scripts/utils": 34
	}],
	2: [function (e, t, n) {
		"use strict";

		function r(e, t) {
			for ("string" != typeof e && (e = e.toString()); e.length < t;) e = "0" + e;
			return e
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var i = e("../../scripts/utils.js"),
			a = function (e) {
				return e && e.__esModule ? e : {
					default: e
				}
			}(i),
			s = e("../../scripts/exceptions"),
			o = {
				year: {
					rx: /Y{1,4}/,
					fn: function (e, t) {
						for (var n = e.getFullYear().toString(); n.length > t.length;) n = n.substr(1, n.length);
						return n
					}
				},
				month: {
					rx: /M{1,4}/,
					fn: function (e, t, n, i) {
						var a = (e.getMonth() + 1).toString();
						switch (t.length) {
							case 1:
								return a;
							case 2:
								return r(a, 2);
							case 3:
								return a = e.getMonth(), i.monthShort[a];
							case 4:
								return a = e.getMonth(), i.month[a]
						}
					}
				},
				day: {
					rx: /D{1,4}/,
					fn: function (e, t, n, i) {
						var a = e.getDate().toString();
						switch (t.length) {
							case 1:
								return a;
							case 2:
								return r(a.toString(), 2);
							case 3:
								return a = e.getDay(), i.weekShort[a];
							case 4:
								return a = e.getDay(), i.week[a]
						}
					}
				},
				hour: {
					rx: /h{1,2}/i,
					fn: function (e, t, n) {
						var r = e.getHours(),
							i = /t{1,2}/i.test(n);
						for (i && r > 12 && (r %= 12), r = r.toString(); r.length < t.length;) r = "0" + r;
						return r
					}
				},
				minute: {
					rx: /m{1,2}/,
					fn: function (e, t) {
						for (var n = e.getMinutes().toString(); n.length < t.length;) n = "0" + n;
						return n
					}
				},
				second: {
					rx: /s{1,2}/,
					fn: function (e, t) {
						for (var n = e.getSeconds().toString(); n.length < t.length;) n = "0" + n;
						return n
					}
				},
				millisecond: {
					rx: /f{1,4}/,
					fn: function (e, t) {
						for (var n = t.length, r = e.getMilliseconds().toString(); r.length < n;) r = "0" + r;
						return r.length > n ? r.substr(0, n) : r
					}
				},
				hoursoffset: {
					rx: /z{1,3}/i,
					fn: function (e, t, n) {
						var i = -e.getTimezoneOffset() / 60,
							a = i > 0 ? "+" : "";
						switch (t.length) {
							case 1:
								return a + i;
							case 2:
								return a + r(i, 2);
							case 3:
								return a + r(i, 2) + ":00"
						}
					}
				},
				ampm: {
					rx: /t{1,2}/i,
					fn: function (e, t) {
						var n, r = e.getHours(),
							i = /T{1,2}/.test(t);
						switch (t.length) {
							case 1:
								n = r > 12 ? "p" : "a";
								break;
							case 2:
								n = r > 12 ? "pm" : "am"
						}
						return i ? n.toUpperCase() : n
					}
				},
				weekday: {
					rx: /w{1,2}/i,
					fn: function (e, t, n, r) {
						var i = e.getDay(),
							a = t.length > 1 ? "week" : "weekShort",
							s = r[a];
						return s && void 0 !== s[i] ? s[i] : i
					}
				}
			},
			u = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+Z?$|^\d{4}-\d{2}-\d{2}[T\s]\d{2}:\d{2}:\d{2}(?:\sUTC)?$/,
			l = /^(\d{4})\D(\d{1,2})\D(\d{1,2})(?:\s(\d{1,2})(?:\D(\d{1,2}))?(?:\D(\d{1,2}))?)?$/;
		n.default = {
			looksLikeDate: function (e) {
				return !!e && (e instanceof Date || "string" == typeof e && (!!l.exec(e) || !!u.exec(e)))
			},
			defaults: {
				format: {
					short: "DD.MM.YYYY",
					long: "DD.MM.YYYY HH:mm:ss"
				},
				week: ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"],
				weekShort: ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"],
				month: ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"],
				monthShort: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"]
			},
			parse: function (e) {
				a.default.isString(e) || (0, s.TypeException)("s", "string");
				var t = l.exec(e);
				if (t) {
					var n = t[4];
					if (n) {
						return new Date(parseInt(t[1]), parseInt(t[2]) - 1, parseInt(t[3]), parseInt(n), parseInt(t[5] || 0), parseInt(t[6] || 0))
					}
					return new Date(t[1], t[2] - 1, t[3])
				}
				if (u.exec(e)) return /Z$/.test(e) || -1 != e.indexOf("UTC") || (e += "Z"), new Date(e)
			},
			format: function (e, t, n) {
				t || (t = this.defaults.format.short), n || (n = this.defaults);
				var r = t;
				for (var i in o) {
					var a = o[i],
						s = t.match(a.rx);
					s && (r = r.replace(a.rx, a.fn(e, s[0], t, n)))
				}
				return r
			},
			formatWithTime: function (e, t) {
				return this.format(e, this.defaults.format.long, t)
			},
			isValid: function (e) {
				return e instanceof Date && isFinite(e)
			},
			sameDay: function (e, t) {
				return e.getFullYear() === t.getFullYear() && e.getMonth() === t.getMonth() && e.getDate() === t.getDate()
			},
			isToday: function (e) {
				return this.sameDay(e, new Date)
			},
			hasTime: function (e) {
				var t = e.getHours(),
					n = e.getMinutes(),
					r = e.getSeconds();
				return !!(t || n || r)
			},
			toIso8601: function (e) {
				return this.format(e, "YYYY-MM-DD") + "T" + this.format(e, "hh:mm:ss") + "." + this.format(e, "fff") + "Z"
			}
		}
	}, {
		"../../scripts/exceptions": 20,
		"../../scripts/utils.js": 34
	}],
	3: [function (e, t, n) {
		"use strict";

		function r(e, t) {
			if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var i = function () {
				function e(e, t) {
					for (var n = 0; n < t.length; n++) {
						var r = t[n];
						r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
					}
				}
				return function (t, n, r) {
					return n && e(t.prototype, n), r && e(t, r), t
				}
			}(),
			a = "function" == typeof Symbol && "symbol" == typeof Symbol.iterator ? function (e) {
				return typeof e
			} : function (e) {
				return e && "function" == typeof Symbol && e.constructor === Symbol && e !== Symbol.prototype ? "symbol" : typeof e
			},
			s = e("../../scripts/utils"),
			o = function (e) {
				return e && e.__esModule ? e : {
					default: e
				}
			}(s),
			u = [],
			l = u.slice,
			f = /\s+/,
			c = function (e, t, n, r) {
				if (!n) return !0;
				if ("object" === (void 0 === n ? "undefined" : a(n))) {
					for (var i in n) e[t].apply(e, [i, n[i]].concat(r));
					return !1
				}
				if (f.test(n)) {
					for (var s = n.split(f), o = 0, u = s.length; o < u; o++) e[t].apply(e, [s[o]].concat(r));
					return !1
				}
				return !0
			},
			d = function (e, t) {
				var n, r = -1,
					i = e.length,
					a = t[0],
					s = t[1],
					o = t[2];
				switch (t.length) {
					case 0:
						for (; ++r < i;)(n = e[r]).callback.call(n.ctx);
						return;
					case 1:
						for (; ++r < i;)(n = e[r]).callback.call(n.ctx, a);
						return;
					case 2:
						for (; ++r < i;)(n = e[r]).callback.call(n.ctx, a, s);
						return;
					case 3:
						for (; ++r < i;)(n = e[r]).callback.call(n.ctx, a, s, o);
						return;
					default:
						for (; ++r < i;)(n = e[r]).callback.apply(n.ctx, t)
				}
			},
			p = function () {
				function e() {
					r(this, e)
				}
				return i(e, [{
					key: "on",
					value: function (e, t, n) {
						return c(this, "on", e, [t, n]) && t ? (this._events || (this._events = {}), (this._events[e] || (this._events[e] = [])).push({
							callback: t,
							context: n,
							ctx: n || this
						}), this) : this
					}
				}, {
					key: "once",
					value: function (e, t, n) {
						if (!c(this, "once", e, [t, n]) || !t) return this;
						var r = this,
							i = o.default.once(function () {
								r.off(e, i), t.apply(this, arguments)
							});
						return i._callback = t, this.on(e, i, n)
					}
				}, {
					key: "off",
					value: function (e, t, n) {
						var r, i, a, s, u, l, f, d;
						if (!this._events || !c(this, "off", e, [t, n])) return this;
						if (!e && !t && !n) return this._events = {}, this;
						for (s = e ? [e] : o.default.keys(this._events), u = 0, l = s.length; u < l; u++)
							if (e = s[u], a = this._events[e]) {
								if (this._events[e] = r = [], t || n)
									for (f = 0, d = a.length; f < d; f++) i = a[f], (t && t !== i.callback && t !== i.callback._callback || n && n !== i.context) && r.push(i);
								r.length || delete this._events[e]
							}
						return this
					}
				}, {
					key: "trigger",
					value: function (e) {
						if (!this._events) return this;
						var t = l.call(arguments, 1);
						if (!c(this, "trigger", e, t)) return this;
						var n = this._events[e],
							r = this._events.all;
						return n && d(n, t), r && d(r, arguments), this
					}
				}, {
					key: "emit",
					value: function (e) {
						return this.trigger(e)
					}
				}, {
					key: "stopListening",
					value: function (e, t, n) {
						var r = this._listeners;
						if (!r) return this;
						var i = !t && !n;
						"object" === (void 0 === t ? "undefined" : a(t)) && (n = this), e && ((r = {})[e._listenerId] = e);
						for (var s in r) r[s].off(t, n, this), i && delete this._listeners[s];
						return this
					}
				}, {
					key: "listenTo",
					value: function (e, t, n) {
						if (2 == arguments.length && "object" == (void 0 === t ? "undefined" : a(t))) {
							var r;
							for (r in t) this.listenTo(e, r, t[r]);
							return this
						}
						return (this._listeners || (this._listeners = {}))[e._listenerId || (e._listenerId = o.default.uniqueId("l"))] = e, "object" === (void 0 === t ? "undefined" : a(t)) && (n = this), e.on(t, n, this), this
					}
				}, {
					key: "listenToOnce",
					value: function (e, t, n) {
						return (this._listeners || (this._listeners = {}))[e._listenerId || (e._listenerId = o.default.uniqueId("l"))] = e, "object" === (void 0 === t ? "undefined" : a(t)) && (n = this), e.once(t, n, this), this
					}
				}]), e
			}();
		n.default = p
	}, {
		"../../scripts/utils": 34
	}],
	4: [function (e, t, n) {
		"use strict";
		Object.defineProperty(n, "__esModule", {
			value: !0
		}), n.default = {
			format: function (e, t) {
				return t || (t = {}), "undefined" != typeof Intl ? Intl.NumberFormat(t.locale || "en-GB").format(e) : (e || "").toString()
			}
		}
	}, {}],
	5: [function (e, t, n) {
		"use strict";
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var r = e("../../scripts/utils"),
			i = function (e) {
				return e && e.__esModule ? e : {
					default: e
				}
			}(r);
		n.default = {
			getPropertyValue: function (e, t) {
				for (var n, r = t.split("."), a = e;
					(n = r.shift()) && (i.default.has(a, n) && (a = a[n]), !i.default.isArray(a)););
				return i.default.isArray(a) && r.length ? this.getCollectionPropertiesValue(a, r.join(".")) : a
			},
			getCollectionPropertiesValue: function (e, t, n) {
				if (!t) return e;
				"boolean" != typeof n && (n = !1);
				for (var r = t.split("."), a = [], s = 0, o = e.length; s < o; s++) {
					var u = e[s];
					if (i.default.has(u, r[0]))
						if (i.default.isArray(u)) {
							var l = this.getCollectionPropertiesValue(u, t);
							(n || l.length) && a.push(l)
						} else if (i.default.isPlainObject(u)) {
						var f = this.getPropertyValue(u, t);
						(n || this.validateValue(f)) && a.push(f)
					} else(n || this.validateValue(u)) && a.push(u);
					else n && a.push(null)
				}
				return a
			},
			validateValue: function (e) {
				return !!e && (!i.default.isArray(e) || !!e.length)
			}
		}
	}, {
		"../../scripts/utils": 34
	}],
	6: [function (e, t, n) {
		"use strict";
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var r = e("../../scripts/utils.js"),
			i = function (e) {
				return e && e.__esModule ? e : {
					default: e
				}
			}(r);
		n.default = {
			getPatternFromStrings: function (e) {
				var t = this;
				if (!e || !e.length) throw new Error("invalid parameter");
				var n = i.default.map(e, function (e) {
					return t.escapeCharsForRegex(e)
				}).join("|");
				return new RegExp("(" + n + ")", "mgi")
			},
			escapeCharsForRegex: function (e) {
				return "string" != typeof e ? "" : e.replace(/([\^\$\.\(\)\[\]\?\!\*\+\{\}\|\/\\])/g, "\\$1").replace(/\s/g, "\\s")
			},
			getSearchPattern: function (e, t) {
				if (!e) return /.+/gim;
				switch (t = i.default.extend({
					searchMode: "fullstring"
				}, t || {}), t.searchMode.toLowerCase()) {
					case "splitwords":
						throw new Error("Not implemented");
					case "fullstring":
						e = this.escapeCharsForRegex(e);
						try {
							return new RegExp("(" + e + ")", "mgi")
						} catch (e) {
							return
						}
						break;
					default:
						throw "invalid searchMode"
				}
			},
			getMatchPattern: function (e) {
				return e ? (e = this.escapeCharsForRegex(e), new RegExp(e, "i")) : /.+/gm
			}
		}
	}, {
		"../../scripts/utils.js": 34
	}],
	7: [function (e, t, n) {
		"use strict";

		function r(e) {
			return e && e.__esModule ? e : {
				default: e
			}
		}

		function i(e) {
			return e.toLowerCase()
		}

		function a(e) {
			return e.toUpperCase()
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var s = e("../../scripts/exceptions"),
			o = e("../../scripts/utils.js"),
			u = r(o),
			l = e("../../scripts/components/string.normalize"),
			f = r(l),
			c = "replace",
			d = "invalid filler (must be as single character)";
		u.default.isString;
		n.default = {
			normalize: f.default,
			replaceAt: function (e, t, n) {
				return e ? e.substr(0, t) + n + e.substr(t + n.length) : e
			},
			findDiacritics: function (e) {
				if (!e) return e;
				for (var t, n = /[^\u0000-\u007E]/gm, r = []; t = n.exec(e);) r.push({
					i: t.index,
					v: t[0]
				});
				return r
			},
			restoreDiacritics: function (e, t, n) {
				if (!e) return e;
				var r = t.length;
				if (!r) return e;
				void 0 === n && (n = 0);
				for (var i, a = n + e.length - 1, s = 0; s < r && (i = t[s], !(i.i > a)); s++) e = this.replaceAt(e, i.i - n, i.v);
				return e
			},
			snakeCase: function (e) {
				return e ? this.removeMultipleSpaces(e.trim())[c](/[^a-zA-Z0-9]/g, "_")[c](/([a-z])[\s\-]?([A-Z])/g, function (e, t, n) {
					return t + "_" + i(n)
				})[c](/([A-Z]+)/g, function (e, t) {
					return i(t)
				})[c](/_{2,}/g, "_") : e
			},
			kebabCase: function (e) {
				return e ? this.removeMultipleSpaces(e.trim())[c](/[^a-zA-Z0-9]/g, "-")[c](/([a-z])[\s\-]?([A-Z])/g, function (e, t, n) {
					return t + "-" + i(n)
				})[c](/([A-Z]+)/g, function (e, t) {
					return i(t)
				})[c](/-{2,}/g, "-") : ""
			},
			camelCase: function (e) {
				return e ? this.removeMultipleSpaces(e.trim())[c](/[^a-zA-Z0-9]+([a-zA-Z])?/g, function (e, t) {
					return a(t)
				})[c](/([a-z])[\s\-]?([A-Z])/g, function (e, t, n) {
					return t + a(n)
				})[c](/^([A-Z]+)/g, function (e, t) {
					return i(t)
				}) : e
			},
			format: function (e) {
				var t = Array.prototype.slice.call(arguments, 1);
				return e[c](/{(\d+)}/g, function (e, n) {
					return void 0 !== t[n] ? t[n] : e
				})
			},
			getString: function (e) {
				return "string" == typeof e ? e : e.toString ? e.toString() : ""
			},
			compare: function (e, t, n, r) {
				n = u.default.isNumber(n) ? n : /^asc/i.test(n) ? 1 : -1;
				var i = u.default.extend({
					ci: !0
				}, r);
				return e && !t ? n : !e && t ? -n : e || t ? e == t ? 0 : (u.default.isString(e) || (e = e.toString()), u.default.isString(t) || (t = t.toString()), i.ci && (e = e.toLowerCase(), t = t.toLowerCase()), (0, f.default)(e) < (0, f.default)(t) ? -n : n) : 0
			},
			ofLength: function (e, t) {
				return new Array(t + 1).join(e)
			},
			center: function (e, t, n) {
				if (t <= 0) throw new Error("length must be > 0");
				if (n || (n = " "), !e) return this.ofLength(n, t);
				if (1 != n.length) throw new Error(d);
				for (var r = Math.floor((t - e.length) / 2), i = this.ofLength(n, r), a = !1, s = i + e + i; s.length < t;) a ? s = fillter + s : s += n, a = !a;
				return s
			},
			startsWith: function (e, t, n) {
				return !(!e || !t) && (n ? 0 == e.toLowerCase().indexOf(t) : 0 == e.indexOf(t))
			},
			ljust: function (e, t, n) {
				if (t <= 0) throw new Error("length must be > 0");
				if (n || (n = " "), !e) return this.ofLength(n, t);
				if (1 != n.length) throw new Error(d);
				for (; e.length < t;) e += n;
				return e
			},
			rjust: function (e, t, n) {
				if (t <= 0) throw new Error("length must be > 0");
				if (n || (n = " "), !e) return this.ofLength(n, t);
				if (1 != n.length) throw new Error(d);
				for (; e.length < t;) e = n + e;
				return e
			},
			removeMultipleSpaces: function (e) {
				return e[c](/\s{2,}/g, " ")
			},
			removeLeadingSpaces: function (e) {
				return e[c](/^\s+|\s+$/, "")
			},
			fixWidth: function (e, t) {
				if (!e) return e;
				t || (t = " ");
				var n, r;
				u.default.isString(e) ? (n = e.split(/\n/g), r = !0) : u.default.isArray(e) ? (n = u.default.clone(e), r = !1) : (0, s.ArgumentException)("s", "expected string or string[]");
				for (var i, a = n.length, o = u.default.max(n, function (e) {
						return e.length
					}), l = 0; l < a; l++) {
					for (i = n[l]; i.length < o;) i += t;
					n[l] = i
				}
				return r ? n.join("\n") : n
			},
			linesWidths: function (e) {
				if (!e) return 0;
				var t;
				return u.default.isString(e) ? t = e.split(/\n/g) : u.default.isArray(e) ? t = u.default.clone(e) : (0, s.ArgumentException)("s", "expected string or string[]"), u.default.map(t, function (e) {
					return e.length
				})
			}
		}
	}, {
		"../../scripts/components/string.normalize": 8,
		"../../scripts/exceptions": 20,
		"../../scripts/utils.js": 34
	}],
	8: [function (e, t, n) {
		"use strict";

		function r(e) {
			return e.replace(/[^\u0000-\u007E]/g, function (e) {
				return a[e] || e
			})
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		for (var i = [{
				base: "A",
				letters: "AⒶＡÀÁÂẦẤẪẨÃĀĂẰẮẴẲȦǠÄǞẢÅǺǍȀȂẠẬẶḀĄȺⱯ"
			}, {
				base: "AA",
				letters: "Ꜳ"
			}, {
				base: "AE",
				letters: "ÆǼǢ"
			}, {
				base: "AO",
				letters: "Ꜵ"
			}, {
				base: "AU",
				letters: "Ꜷ"
			}, {
				base: "AV",
				letters: "ꜸꜺ"
			}, {
				base: "AY",
				letters: "Ꜽ"
			}, {
				base: "B",
				letters: "BⒷＢḂḄḆɃƂƁ"
			}, {
				base: "C",
				letters: "CⒸＣĆĈĊČÇḈƇȻꜾ"
			}, {
				base: "D",
				letters: "DⒹＤḊĎḌḐḒḎĐƋƊƉꝹÐ"
			}, {
				base: "DZ",
				letters: "ǱǄ"
			}, {
				base: "Dz",
				letters: "ǲǅ"
			}, {
				base: "E",
				letters: "EⒺＥÈÉÊỀẾỄỂẼĒḔḖĔĖËẺĚȄȆẸỆȨḜĘḘḚƐƎ"
			}, {
				base: "F",
				letters: "FⒻＦḞƑꝻ"
			}, {
				base: "G",
				letters: "GⒼＧǴĜḠĞĠǦĢǤƓꞠꝽꝾ"
			}, {
				base: "H",
				letters: "HⒽＨĤḢḦȞḤḨḪĦⱧⱵꞍ"
			}, {
				base: "I",
				letters: "IⒾＩÌÍÎĨĪĬİÏḮỈǏȈȊỊĮḬƗ"
			}, {
				base: "J",
				letters: "JⒿＪĴɈ"
			}, {
				base: "K",
				letters: "KⓀＫḰǨḲĶḴƘⱩꝀꝂꝄꞢ"
			}, {
				base: "L",
				letters: "LⓁＬĿĹĽḶḸĻḼḺŁȽⱢⱠꝈꝆꞀ"
			}, {
				base: "LJ",
				letters: "Ǉ"
			}, {
				base: "Lj",
				letters: "ǈ"
			}, {
				base: "M",
				letters: "MⓂＭḾṀṂⱮƜ"
			}, {
				base: "N",
				letters: "NⓃＮǸŃÑṄŇṆŅṊṈȠƝꞐꞤ"
			}, {
				base: "NJ",
				letters: "Ǌ"
			}, {
				base: "Nj",
				letters: "ǋ"
			}, {
				base: "O",
				letters: "OⓄＯÒÓÔỒỐỖỔÕṌȬṎŌṐṒŎȮȰÖȪỎŐǑȌȎƠỜỚỠỞỢỌỘǪǬØǾƆƟꝊꝌ"
			}, {
				base: "OI",
				letters: "Ƣ"
			}, {
				base: "OO",
				letters: "Ꝏ"
			}, {
				base: "OU",
				letters: "Ȣ"
			}, {
				base: "OE",
				letters: "Œ"
			}, {
				base: "oe",
				letters: "œ"
			}, {
				base: "P",
				letters: "PⓅＰṔṖƤⱣꝐꝒꝔ"
			}, {
				base: "Q",
				letters: "QⓆＱꝖꝘɊ"
			}, {
				base: "R",
				letters: "RⓇＲŔṘŘȐȒṚṜŖṞɌⱤꝚꞦꞂ"
			}, {
				base: "S",
				letters: "SⓈＳẞŚṤŜṠŠṦṢṨȘŞⱾꞨꞄ"
			}, {
				base: "T",
				letters: "TⓉＴṪŤṬȚŢṰṮŦƬƮȾꞆ"
			}, {
				base: "TZ",
				letters: "Ꜩ"
			}, {
				base: "U",
				letters: "UⓊＵÙÚÛŨṸŪṺŬÜǛǗǕǙỦŮŰǓȔȖƯỪỨỮỬỰỤṲŲṶṴɄ"
			}, {
				base: "V",
				letters: "VⓋＶṼṾƲꝞɅ"
			}, {
				base: "VY",
				letters: "Ꝡ"
			}, {
				base: "W",
				letters: "WⓌＷẀẂŴẆẄẈⱲ"
			}, {
				base: "X",
				letters: "XⓍＸẊẌ"
			}, {
				base: "Y",
				letters: "YⓎＹỲÝŶỸȲẎŸỶỴƳɎỾ"
			}, {
				base: "Z",
				letters: "ZⓏＺŹẐŻŽẒẔƵȤⱿⱫꝢ"
			}, {
				base: "a",
				letters: "aⓐａẚàáâầấẫẩãāăằắẵẳȧǡäǟảåǻǎȁȃạậặḁąⱥɐ"
			}, {
				base: "aa",
				letters: "ꜳ"
			}, {
				base: "ae",
				letters: "æǽǣ"
			}, {
				base: "ao",
				letters: "ꜵ"
			}, {
				base: "au",
				letters: "ꜷ"
			}, {
				base: "av",
				letters: "ꜹꜻ"
			}, {
				base: "ay",
				letters: "ꜽ"
			}, {
				base: "b",
				letters: "bⓑｂḃḅḇƀƃɓ"
			}, {
				base: "c",
				letters: "cⓒｃćĉċčçḉƈȼꜿↄ"
			}, {
				base: "d",
				letters: "dⓓｄḋďḍḑḓḏđƌɖɗꝺ"
			}, {
				base: "dz",
				letters: "ǳǆ"
			}, {
				base: "e",
				letters: "eⓔｅèéêềếễểẽēḕḗĕėëẻěȅȇẹệȩḝęḙḛɇɛǝ"
			}, {
				base: "f",
				letters: "fⓕｆḟƒꝼ"
			}, {
				base: "g",
				letters: "gⓖｇǵĝḡğġǧģǥɠꞡᵹꝿ"
			}, {
				base: "h",
				letters: "hⓗｈĥḣḧȟḥḩḫẖħⱨⱶɥ"
			}, {
				base: "hv",
				letters: "ƕ"
			}, {
				base: "i",
				letters: "iⓘｉìíîĩīĭïḯỉǐȉȋịįḭɨı"
			}, {
				base: "j",
				letters: "jⓙｊĵǰɉ"
			}, {
				base: "k",
				letters: "kⓚｋḱǩḳķḵƙⱪꝁꝃꝅꞣ"
			}, {
				base: "l",
				letters: "lⓛｌŀĺľḷḹļḽḻſłƚɫⱡꝉꞁꝇ"
			}, {
				base: "lj",
				letters: "ǉ"
			}, {
				base: "m",
				letters: "mⓜｍḿṁṃɱɯ"
			}, {
				base: "n",
				letters: "nⓝｎǹńñṅňṇņṋṉƞɲŉꞑꞥ"
			}, {
				base: "nj",
				letters: "ǌ"
			}, {
				base: "o",
				letters: "oⓞｏòóôồốỗổõṍȭṏōṑṓŏȯȱöȫỏőǒȍȏơờớỡởợọộǫǭøǿɔꝋꝍɵ"
			}, {
				base: "oi",
				letters: "ƣ"
			}, {
				base: "ou",
				letters: "ȣ"
			}, {
				base: "oo",
				letters: "ꝏ"
			}, {
				base: "p",
				letters: "pⓟｐṕṗƥᵽꝑꝓꝕ"
			}, {
				base: "q",
				letters: "qⓠｑɋꝗꝙ"
			}, {
				base: "r",
				letters: "rⓡｒŕṙřȑȓṛṝŗṟɍɽꝛꞧꞃ"
			}, {
				base: "s",
				letters: "sⓢｓßśṥŝṡšṧṣṩșşȿꞩꞅẛ"
			}, {
				base: "t",
				letters: "tⓣｔṫẗťṭțţṱṯŧƭʈⱦꞇ"
			}, {
				base: "tz",
				letters: "ꜩ"
			}, {
				base: "u",
				letters: "uⓤｕùúûũṹūṻŭüǜǘǖǚủůűǔȕȗưừứữửựụṳųṷṵʉ"
			}, {
				base: "v",
				letters: "vⓥｖṽṿʋꝟʌ"
			}, {
				base: "vy",
				letters: "ꝡ"
			}, {
				base: "w",
				letters: "wⓦｗẁẃŵẇẅẘẉⱳ"
			}, {
				base: "x",
				letters: "xⓧｘẋẍ"
			}, {
				base: "y",
				letters: "yⓨｙỳýŷỹȳẏÿỷẙỵƴɏỿ"
			}, {
				base: "z",
				letters: "zⓩｚźẑżžẓẕƶȥɀⱬꝣ"
			}], a = {}, s = 0; s < i.length; s++)
			for (var o = i[s].letters, u = 0; u < o.length; u++) a[o[u]] = i[s].base;
		n.default = r
	}, {}],
	9: [function (e, t, n) {
		"use strict";

		function r(e) {
			return e && e.__esModule ? e : {
				default: e
			}
		}

		function i(e) {
			return e.indexOf("json") > -1 && e != l ? l : e
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var a = e("../../scripts/utils"),
			s = r(a),
			o = e("../../scripts/data/json"),
			u = r(o),
			l = "application/json",
			f = {
				type: "POST",
				headers: {
					"X-Requested-With": "XMLHttpRequest",
					"Content-Type": l
				},
				json: {
					parseDates: !0
				}
			};
		n.default = {
			defaults: f,
			requestBeforeSend: function (e, t, n) {},
			setup: function (e) {
				if (!s.default.isPlainObject(e)) throw new Error("Invalid options for AJAX setup.");
				return s.default.extend(this.defaults, e), this
			},
			converters: {
				"application/json": function (e, t, n) {
					return u.default.parse(e, n.json)
				}
			},
			createQs: function (e) {
				if (!e) return "";
				var t, n, r = [];
				for (t in e) n = e[t], s.default.isNullOrEmptyString(n) || r.push([t, n]);
				return r.sort(function (e, t) {
					return e > t ? 1 : e < t ? -1 : 0
				}), s.default.map(r, function (e) {
					return encodeURIComponent(e[0]) + "=" + encodeURIComponent(e[1])
				}).join("&")
			},
			shot: function (e) {
				if (e || (e = {}), e.headers) {
					e.headers
				}
				var t = s.default.extend({}, s.default.clone(this.defaults), e);
				e.headers && (t.headers = s.default.extend({}, this.defaults.headers, e.headers));
				var n = t.url;
				if (!n) throw new Error("missing `url` for XMLHttpRequest");
				var r = this,
					a = r.converters,
					o = t.type;
				if (!o) throw new Error("missing `type` for XMLHttpRequest");
				var l = "GET" == o,
					f = t.data;
				if (l && f) {
					var c = this.createQs(f),
						d = -1 != n.indexOf("?");
					n += (d ? "&" : "?") + c, delete t.headers["Content-Type"]
				}
				return new Promise(function (f, c) {
					var d = new XMLHttpRequest;
					d.open(o, n);
					var p = t.headers;
					if (p) {
						var h;
						for (h in p) d.setRequestHeader(h, p[h])
					}
					d.onload = function () {
						if (200 == d.status) {
							var e = d.response,
								n = i(d.getResponseHeader("Content-Type") || ""),
								r = a[n];
							s.default.isFunction(r) && (e = r(e, d, t)), f(e, d.status, d)
						} else c(Error(d.statusText))
					}, d.onerror = function () {
						c(d, null, Error("Network Error"))
					};
					var m = t.data;
					if (m && !l) {
						var g = t.headers["Content-Type"];
						if (!(g.indexOf("/json") > -1)) throw "application/x-www-form-urlencoded; charset=UTF-8" == g ? "Not implemented" : "invalid or not implemented content type: " + g;
						m = u.default.compose(m), r.requestBeforeSend(d, t, e), d.send(m)
					} else r.requestBeforeSend(d, t, e), d.send()
				})
			},
			get: function (e, t) {
				return t = t || {}, t.url = e, t.type = "GET", this.shot(t)
			},
			post: function (e, t) {
				return t = t || {}, t.url = e, t.type = "POST", this.shot(t)
			}
		}
	}, {
		"../../scripts/data/json": 13,
		"../../scripts/utils": 34
	}],
	10: [function (e, t, n) {
		"use strict";
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var r = e("../../scripts/utils"),
			i = function (e) {
				return e && e.__esModule ? e : {
					default: e
				}
			}(r),
			a = {
				allStrings: 1,
				keepType: 2
			};
		n.default = {
			default: {
				addBom: !0,
				separator: ",",
				addSeparatorLine: !1,
				typeHandling: a.keepType
			},
			serialize: function (e, t) {
				for (var n = i.default.extend({}, this.default, t), r = [], s = n.separator, o = n.typeHandling, u = n.addBom ? "\ufeff" : "", l = 0, f = e.length; l < f; l++) {
					for (var c = [], d = e[l], p = 0, h = d.length; p < h; p++) {
						var m = d[p];
						m instanceof Date ? m = m.toLocaleString() : "string" != typeof m && (m = m && m.toString ? m.toString() : ""), /"/.test(m) && (m = m.replace(/"/g, '""')), (o == a.allStrings || /"|\n/.test(m) || m.indexOf(s) > -1) && (m = '"' + m + '"'), c.push(m)
					}
					r.push(c.join(s))
				}
				return n.addSeparatorLine && r.push("\t" + s), u + r.join("\n")
			}
		}
	}, {
		"../../scripts/utils": 34
	}],
	11: [function (e, t, n) {
		"use strict";
		Object.defineProperty(n, "__esModule", {
			value: !0
		}), n.default = {
			supportsCsExport: function () {
				return navigator.msSaveBlob || function () {
					return void 0 !== document.createElement("a").download
				}()
			},
			exportfile: function (e, t, n) {
				var r = new Blob([t], {
					type: n
				});
				if (navigator.msSaveBlob) navigator.msSaveBlob(r, e);
				else {
					var i = document.createElement("a");
					if (void 0 !== i.download) {
						var a = URL.createObjectURL(r);
						i.setAttribute("href", a), i.setAttribute("download", e);
						var s = {
							visibility: "hidden",
							position: "absolute",
							left: "-9999px"
						};
						for (var o in s) i.style[o] = s[o];
						document.body.appendChild(i), i.click(), document.body.removeChild(i)
					}
				}
			}
		}
	}, {}],
	12: [function (e, t, n) {
		"use strict";

		function r(e, t) {
			if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
			return !t || "object" != typeof t && "function" != typeof t ? e : t
		}

		function i(e, t) {
			if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
			e.prototype = Object.create(t && t.prototype, {
				constructor: {
					value: e,
					enumerable: !1,
					writable: !0,
					configurable: !0
				}
			}), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
		}

		function a(e, t) {
			if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
		}

		function s(e) {
			return "number" == typeof e
		}

		function o(e) {
			return e.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;").replace(/'/g, "&#039;")
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var u = function () {
				function e(e, t) {
					for (var n = 0; n < t.length; n++) {
						var r = t[n];
						r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
					}
				}
				return function (t, n, r) {
					return n && e(t.prototype, n), r && e(t, r), t
				}
			}(),
			l = function () {
				function e(t, n, r) {
					a(this, e), this.tagName = t, this.attributes = n || {}, r && r instanceof Array == 0 && (r = [r]), this.children = r || [], this.hidden = !1, this.empty = !1
				}
				return u(e, [{
					key: "appendChild",
					value: function (e) {
						this.children.push(e)
					}
				}, {
					key: "toString",
					value: function (e, t, n) {
						var r, i = this,
							a = i.empty,
							o = i.tagName,
							u = i.attributes,
							e = s(e) ? e : 0,
							n = s(n) ? n : 0,
							l = e > 0 ? new Array(e * n + 1).join(t || " ") : "",
							f = "<" + o;
						for (r in u) void 0 !== u[r] && (c.indexOf(r) > -1 ? f += " " + r : f += " " + r + '="' + u[r] + '"');
						if (a) return f += " />", e > 0 && (f = l + f + "\n"), f;
						f += ">";
						var d = i.children;
						e > 0 && d.length && (f += "\n");
						for (var p = 0, h = d.length; p < h; p++) {
							var m = d[p];
							m && ("string" == typeof m ? f += l + m + "\n" : m.hidden || (f += m.toString(e, t, n + 1)))
						}
						return d && d.length ? f += l + "</" + o + ">" : f += "</" + o + ">", e > 0 && (f = l + f + "\n"), f
					}
				}, {
					key: "tagName",
					get: function () {
						return this._tagName
					},
					set: function (e) {
						if ("string" != typeof e) throw new Error("tagName must be a string");
						if (!e.trim()) throw new Error("tagName must have a length");
						if (e.indexOf(" ") > -1) throw new Error("tagName cannot contain spaces");
						this._tagName = e
					}
				}]), e
			}(),
			f = "area base basefont br col frame hr img input isindex link meta param".split(" "),
			c = "checked selected disabled readonly multiple ismap isMap defer noresize noResize nowrap noWrap noshade noShade compact".split(" "),
			d = function (e) {
				function t(e, n, i) {
					a(this, t);
					var s = r(this, (t.__proto__ || Object.getPrototypeOf(t)).call(this, e, n, i));
					return s.empty = f.indexOf(e.toLowerCase()) > -1, s
				}
				return i(t, e), u(t, [{
					key: "id",
					get: function () {
						return this.attributes.id
					},
					set: function (e) {
						this.attributes.id = e
					}
				}]), t
			}(l),
			p = function () {
				function e(t) {
					a(this, e), this.text = t
				}
				return u(e, [{
					key: "toString",
					value: function (e, t, n) {
						var e = s(e) ? e : 0,
							n = s(n) ? n : 0,
							r = e > 0 ? new Array(e * n + 1).join(t || " ") : "";
						return r + this.text + (r ? "\n" : "")
					}
				}, {
					key: "text",
					get: function () {
						return this._text
					},
					set: function (e) {
						e || (e = ""), "string" != typeof e && (e = e.toString()), this._text = o(e)
					}
				}]), e
			}(),
			h = function () {
				function e(t) {
					a(this, e), this.html = t
				}
				return u(e, [{
					key: "toString",
					value: function (e, t, n) {
						var e = s(e) ? e : 0,
							n = s(n) ? n : 0,
							r = e > 0 ? new Array(e * n + 1).join(t || " ") : "";
						return r + this.html + (r ? "\n" : "")
					}
				}]), e
			}(),
			m = function () {
				function e(t) {
					a(this, e), this.text = t
				}
				return u(e, [{
					key: "toString",
					value: function (e, t, n) {
						var e = s(e) ? e : 0,
							n = s(n) ? n : 0,
							r = e > 0 ? new Array(e * n + 1).join(t || " ") : "";
						return r + "\x3c!--" + this.text + "--\x3e" + (r ? "\n" : "")
					}
				}, {
					key: "text",
					get: function () {
						return this._text
					},
					set: function (e) {
						e || (e = ""), "string" != typeof e && (e = e.toString()), e = e.replace(/<!--/g, "").replace(/-->/g, ""), this._text = e
					}
				}]), e
			}(),
			g = function () {
				function e(t) {
					a(this, e), this.children = t, this.hidden = !1
				}
				return u(e, [{
					key: "toString",
					value: function (e, t, n) {
						var r = "",
							i = this.children;
						if (!i || this.hidden) return r;
						for (var a = 0, s = i.length; a < s; a++) {
							var o = i[a];
							o && (o.hidden || (r += o.toString(e, t, n)))
						}
						return r
					}
				}]), e
			}();
		n.VXmlElement = l, n.VHtmlElement = d, n.VHtmlFragment = h, n.VTextElement = p, n.VCommentElement = m, n.VWrapperElement = g, n.escapeHtml = o
	}, {}],
	13: [function (e, t, n) {
		"use strict";

		function r(e) {
			return e && e.__esModule ? e : {
				default: e
			}
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var i = e("../../scripts/utils"),
			a = r(i),
			s = e("../../scripts/components/date"),
			o = r(s);
		n.default = {
			compose: function (e, t) {
				return JSON.stringify(e, function (e, t) {
					return void 0 === t ? null : t
				}, t)
			},
			parse: function (e, t) {
				return a.default.extend({
					parseDates: !0
				}, t).parseDates ? JSON.parse(e, function (e, t) {
					if (a.default.isString(t) && o.default.looksLikeDate(t)) {
						var n = o.default.parse(t);
						if (n && o.default.isValid(n)) return n
					}
					return t
				}) : JSON.parse(e)
			},
			clone: function (e) {
				return this.parse(this.compose(e))
			}
		}
	}, {
		"../../scripts/components/date": 2,
		"../../scripts/utils": 34
	}],
	14: [function (e, t, n) {
		"use strict";

		function r(e) {
			return e && e.__esModule ? e : {
				default: e
			}
		}

		function i(e) {
			if (s.default.isObject(e)) return e;
			switch (e) {
				case 1:
					return localStorage;
				case 2:
				default:
					return sessionStorage
			}
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var a = e("../../scripts/utils"),
			s = r(a),
			o = e("../../scripts/data/json"),
			u = r(o);
		e("../../scripts/exceptions");
		n.default = {
			get: function (e, t, n, r) {
				!0 === t && (r = !0, t = void 0), !0 === n && (r = !0, n = void 0);
				var a, o = i(n),
					l = o.getItem(e);
				if (l) {
					try {
						l = u.default.parse(l)
					} catch (t) {
						return void o.removeItem(e)
					}
					if (!t) return s.default.map(l, function (e) {
						return r ? e : e.data
					});
					var f, c = [],
						d = l.length;
					for (a = 0; a < d; a++) {
						var p = l[a];
						if (p) {
							var h = p.data,
								m = h.expiration;
							s.default.isNumber(m) && m > 0 && (new Date).getTime() > m ? c.push(p) : t(h) && (f = r ? p : h)
						}
					}
					return c.length && this.remove(e, function (e) {
						return c.indexOf(e) > -1
					}), f
				}
			},
			remove: function (e, t, n) {
				var r = i(n);
				if (!t) return void r.removeItem(e);
				var a, s = r.getItem(e);
				if (s) {
					try {
						s = u.default.parse(s)
					} catch (t) {
						return void r.removeItem(e)
					}
					var o = s.length,
						l = [];
					for (a = 0; a < o; a++) {
						var f = s[a];
						if (f) {
							t(f.data) || l.push(f)
						}
					}
					return r.setItem(e, u.default.compose(l))
				}
			},
			set: function (e, t, n, r, a) {
				s.default.isNumber(n) || (n = 10), s.default.isNumber(r) || (r = -1);
				var o = i(a),
					l = (new Date).getTime(),
					f = r > 0 ? l + r : -1,
					c = {
						ts: l,
						expiration: f,
						data: t
					},
					d = o.getItem(e);
				if (d) {
					try {
						d = u.default.parse(d)
					} catch (r) {
						return o.removeItem(e), this.set(e, t, n)
					}
					d.length >= n && d.shift(), d.push(c)
				} else d = [{
					ts: l,
					expiration: f,
					data: t
				}];
				return o.setItem(e, u.default.compose(d))
			}
		}
	}, {
		"../../scripts/data/json": 13,
		"../../scripts/exceptions": 20,
		"../../scripts/utils": 34
	}],
	15: [function (e, t, n) {
		"use strict";
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var r = {};
		n.default = {
			items: function () {
				return r
			},
			length: function () {
				var e, t = 0;
				for (e in r) t++;
				return t
			},
			getItem: function (e) {
				return r[e]
			},
			setItem: function (e, t) {
				r[e] = t
			},
			removeItem: function (e) {
				delete r[e]
			},
			clear: function () {
				var e;
				for (e in r) delete r[e]
			}
		}
	}, {}],
	16: [function (e, t, n) {
		"use strict";

		function r(e, t) {
			if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var i = "function" == typeof Symbol && "symbol" == typeof Symbol.iterator ? function (e) {
				return typeof e
			} : function (e) {
				return e && "function" == typeof Symbol && e.constructor === Symbol && e !== Symbol.prototype ? "symbol" : typeof e
			},
			a = function () {
				function e(e, t) {
					for (var n = 0; n < t.length; n++) {
						var r = t[n];
						r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
					}
				}
				return function (t, n, r) {
					return n && e(t.prototype, n), r && e(t, r), t
				}
			}(),
			s = e("../../scripts/utils"),
			o = function (e) {
				return e && e.__esModule ? e : {
					default: e
				}
			}(s),
			u = function () {
				function e() {
					r(this, e)
				}
				return a(e, [{
					key: "describe",
					value: function (e, t) {
						if (o.default.isArray(e)) return this.describeList(e, t);
						var n, r = {};
						for (n in e) r[n] = this.getType(e[n]);
						return r
					}
				}, {
					key: "describeList",
					value: function (e, t) {
						function n(e) {
							return e
						}
						var r = {};
						t = t || {};
						for (var i = o.default.isNumber(t.limit) ? t.limit : e.length, a = 0; a < i; a++) {
							var s = this.describe(e[a]);
							for (var u in s) o.default.has(r, u) ? void 0 != n(s[u]) && n(r[u]) != n(s[u]) && (n(r[u]) ? r[u] = "string" : r[u] = n(s[u])) : o.default.extend(r, s);
							if (t.lazy && !o.default.any(r, function (e, t) {
									return void 0 === t
								})) break
						}
						return r
					}
				}, {
					key: "getType",
					value: function (e) {
						if (null != e && void 0 != e) return e instanceof Array ? "array" : e instanceof Date ? "date" : e instanceof RegExp ? "regex" : void 0 === e ? "undefined" : i(e)
					}
				}]), e
			}();
		n.default = u
	}, {
		"../../scripts/utils": 34
	}],
	17: [function (e, t, n) {
		"use strict";

		function r(e, t) {
			if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var i = function () {
				function e(e, t) {
					for (var n = 0; n < t.length; n++) {
						var r = t[n];
						r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
					}
				}
				return function (t, n, r) {
					return n && e(t.prototype, n), r && e(t, r), t
				}
			}(),
			a = e("../../scripts/utils"),
			s = function (e) {
				return e && e.__esModule ? e : {
					default: e
				}
			}(a),
			o = function () {
				function e() {
					r(this, e)
				}
				return i(e, [{
					key: "sanitize",
					value: function (e) {
						var t;
						for (t in e)
							if (s.default.isString(e[t])) e[t] = this.escape(e[t]);
							else if (s.default.isObject(e[t]))
							if (s.default.isArray(e[t]))
								for (var n = 0, r = e[t].length; n < r; n++) e[t][n] = this.sanitize(e[t][n]);
							else e[t] = this.sanitize(e[t]);
						return e
					}
				}, {
					key: "escapeHtml",
					value: function (e) {
						return e.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;").replace(/'/g, "&#039;")
					}
				}, {
					key: "escape",
					value: function (e) {
						return e ? this.escapeHtml(e) : ""
					}
				}]), e
			}();
		n.default = o
	}, {
		"../../scripts/utils": 34
	}],
	18: [function (e, t, n) {
		"use strict";
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var r = e("../../scripts/utils");
		! function (e) {
			e && e.__esModule
		}(r);
		n.default = {
			normal: function (e) {
				return '<?xml version="1.0"?>' + e.replace(/\sxmlns="http:\/\/www\.w3\.org\/\d+\/xhtml"/, "")
			},
			pretty: function (e, t) {
				e = this.normal(e), "number" != typeof t && (t = 2);
				var n = [];
				e = e.replace(/(>)(<)(\/*)/g, "$1\r\n$2$3");
				for (var r = 0, i = e.split("\r\n"), a = i.length, s = 0; s < a; s++) {
					var o = i[s],
						u = 0;
					o.match(/.+<\/\w[^>]*>$/) ? u = 0 : o.match(/^<\/\w/) ? 0 != r && (r -= 1) : u = o.match(/^<\w[^>]*[^\/]>.*$/) ? 1 : 0;
					var l = new Array(r * t).join(" ");
					n.push(l + o + "\r\n"), r += u
				}
				return n.join("")
			}
		}
	}, {
		"../../scripts/utils": 34
	}],
	19: [function (e, t, n) {
		"use strict";

		function r(e, t, n) {
			if (t.search(/\s/) > -1) {
				t = t.split(/\s/g);
				for (var i = 0, a = t[F]; i < a; i++) r(e, t[i], n)
			} else(void 0 === t ? "undefined" : C(t)) == M && e.classList[n ? "add" : "remove"](t);
			return e
		}

		function i(e, t) {
			return r(e, t, 1)
		}

		function a(e, t) {
			return r(e, t, 0)
		}

		function s(e, t) {
			return e && e.classList.contains(t)
		}

		function o(e, t) {
			return e.getAttribute(t)
		}

		function u(e) {
			return o(e, "name")
		}

		function l(e) {
			return P(e) && "password" == o(e, "type")
		}

		function f(e, t) {
			if ("checkbox" == e.type) return e.checked = 1 == t || /1|true/.test(t), void e.dispatchEvent(new Event("change"), {
				forced: !0
			});
			e.value != t && (e.value = t, e.dispatchEvent(new Event("change"), {
				forced: !0
			}))
		}

		function c(e) {
			if (/input/i.test(e.tagName)) switch (o(e, "type")) {
				case "radio":
				case "checkbox":
					return e.checked
			}
			return e.value
		}

		function d(e) {
			return e && /^input$/i.test(e.tagName) && /^(radio)$/i.test(e.type)
		}

		function p(e) {
			return e && (/^select$/i.test(e.tagName) || d(e))
		}

		function h(e) {
			return e.nextElementSibling
		}

		function m(e) {
			return e.previousElementSibling
		}

		function g(e, t) {
			return e.querySelectorAll(t)
		}

		function v(e, t) {
			return e.querySelectorAll(t)[0]
		}

		function y(e, t) {
			return e.getElementsByClassName(t)[0]
		}

		function b(e) {
			var t = window.getComputedStyle(e);
			return "none" == t.display || "hidden" == t.visibility
		}

		function w(e) {
			return document.createElement(e)
		}

		function k(e, t) {
			e.parentNode.insertBefore(t, e.nextSibling)
		}

		function x(e, t) {
			e.appendChild(t)
		}

		function E(e) {
			return ("undefined" == typeof HTMLElement ? "undefined" : C(HTMLElement)) === j ? e instanceof HTMLElement : e && (void 0 === e ? "undefined" : C(e)) === j && null !== e && 1 === e.nodeType && C(e.nodeName) === M
		}

		function S(e) {
			return e && E(e) && /input|button|textarea|select/i.test(e.tagName)
		}

		function P(e) {
			return e && E(e) && /input/i.test(e.tagName)
		}

		function T(e) {
			if (!E(e)) throw new Error("expected HTML Element");
			var t = e.parentNode;
			if (!E(t)) throw new Error("expected HTML element with parentNode");
			return t
		}

		function _(e) {
			var t = e.indexOf(I);
			if (t > -1) {
				e.substr(0, t);
				return [e.substr(0, t), e.substr(t + 1)]
			}
			return [e, ""]
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var C = "function" == typeof Symbol && "symbol" == typeof Symbol.iterator ? function (e) {
				return typeof e
			} : function (e) {
				return e && "function" == typeof Symbol && e.constructor === Symbol && e !== Symbol.prototype ? "symbol" : typeof e
			},
			O = e("../scripts/utils.js"),
			V = function (e) {
				return e && e.__esModule ? e : {
					default: e
				}
			}(O),
			j = "object",
			M = "string",
			F = "length",
			D = V.default.any,
			H = V.default.each,
			I = ".",
			N = [];
		n.default = {
			splitNamespace: _,
			empty: function (e) {
				for (; e.hasChildNodes();) this.off(e.lastChild), e.removeChild(e.lastChild)
			},
			remove: function (e) {
				if (e) {
					this.off(e);
					var t = e.parentElement || e.parentNode;
					t && t.removeChild(e)
				}
			},
			closest: function (e, t, n) {
				if (e && t) {
					if (!n && t(e)) return e;
					for (var r = e; r = r.parentElement;)
						if (t(r)) return r
				}
			},
			closestWithTag: function (e, t, n) {
				if (t) return t = t.toUpperCase(), this.closest(e, function (e) {
					return e.tagName == t
				}, n)
			},
			closestWithClass: function (e, t, n) {
				if (t) return this.closest(e, function (e) {
					return s(e, t)
				}, n)
			},
			contains: function (e, t) {
				if (!e || !t) return !1;
				if (!e.hasChildNodes()) return !1;
				for (var n = e.childNodes, r = n.length, i = 0; i < r; i++) {
					var a = n[i];
					if (a === t) return !0;
					if (this.contains(a, t)) return !0
				}
				return !1
			},
			on: function (e, t, n, r) {
				if (!E(e)) throw new Error("argument is not a DOM element.");
				V.default.isFunction(n) && !r && (r = n, n = null);
				var i = this,
					a = _(t),
					s = a[0],
					o = a[1],
					u = function (t) {
						t.target;
						if (!n) {
							var a = r(t, t.detail);
							return !1 === a && t.preventDefault(), !0
						}
						var s = g(e, n);
						if (D(s, function (e) {
								return t.target === e || i.contains(e, t.target)
							})) {
							var a = r(t, t.detail);
							return !1 === a && t.preventDefault(), !0
						}
						return !0
					};
				return N.push({
					type: t,
					ev: s,
					ns: o,
					fn: u,
					el: e
				}), e.addEventListener(s, u, !0), this
			},
			off: function (e, t) {
				if (E(e))
					if (t)
						if (t[0] === I) {
							var n = t.substr(1);
							H(N, function (t) {
								t.el === e && t.ns == n && t.el.removeEventListener(t.ev, t.fn, !0)
							})
						} else {
							var r = _(t),
								i = r[0],
								n = r[1];
							H(N, function (t) {
								t.el !== e || t.ev != i || n && t.ns != n || t.el.removeEventListener(t.ev, t.fn, !0)
							})
						}
				else H(N, function (t) {
					t.el === e && t.el.removeEventListener(t.ev, t.fn, !0)
				})
			},
			offAll: function () {
				var e, t = this;
				return H(N, function (t) {
					e = t.el, e.removeEventListener(t.ev, t.fn, !0)
				}), t
			},
			fire: function (e, t, n) {
				if ("focus" == t) return void e.focus();
				var r;
				window.CustomEvent ? r = new CustomEvent(t, {
					detail: n
				}) : document.createEvent && (r = document.createEvent("CustomEvent"), r.initCustomEvent(t, !0, !0, n)), e.dispatchEvent(r)
			},
			siblings: function (e, t) {
				for (var n = T(e), r = [], i = n[t ? "childNodes" : "children"], a = 0, s = i.length; a < s; a++) {
					var o = i[a];
					o !== e && r.push(o)
				}
				return r
			},
			nextSiblings: function (e, t) {
				for (var n = T(e), r = [], i = n[t ? "childNodes" : "children"], a = !1, s = 0, o = i.length; s < o; s++) {
					var u = i[s];
					u !== e && a ? r.push(u) : a = !0
				}
				return r
			},
			prevSiblings: function (e, t) {
				for (var n = T(e), r = [], i = n[t ? "childNodes" : "children"], a = 0, s = i.length; a < s; a++) {
					var o = i[a];
					if (o === e) break;
					r.push(o)
				}
				return r
			},
			findByClass: function (e, t) {
				return e.getElementsByClassName(t)
			},
			isFocused: function (e) {
				return !!e && e === this.getFocusedElement()
			},
			getFocusedElement: function () {
				return document.querySelector(":focus")
			},
			anyInputFocused: function () {
				var e = this.getFocusedElement();
				return e && /input|select|textarea/i.test(e.tagName)
			},
			prev: m,
			next: h,
			append: x,
			addClass: i,
			removeClass: a,
			modClass: r,
			attr: o,
			hasClass: s,
			after: k,
			createElement: w,
			isElement: E,
			isInput: P,
			isAnyInput: S,
			isSelectable: p,
			isRadioButton: d,
			isPassword: l,
			attrName: u,
			isHidden: b,
			find: g,
			findFirst: v,
			findFirstByClass: y,
			getValue: c,
			setValue: f
		}
	}, {
		"../scripts/utils.js": 34
	}],
	20: [function (e, t, n) {
		"use strict";

		function r(e) {
			return "number" == typeof e
		}

		function i(e) {
			throw new Error("The parameter cannot be null: " + (e || l))
		}

		function a(e) {
			throw new Error("Invalid argument: " + (e || l))
		}

		function s(e, t) {
			throw new Error("Expected parameter: " + (e || l) + " of type: " + (type || l))
		}

		function o(e) {
			throw new Error("Invalid operation: " + e)
		}

		function u(e, t, n) {
			var i = "Out of range. Expected parameter: " + (e || l);
			throw r(n) || 0 !== t ? (r(t) && (i = " >=" + t), r(n) && (i = " <=" + n)) : i = " to be positive.", new Error(i)
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var l = "???";
		n.ArgumentException = a, n.ArgumentNullException = i, n.TypeException = s, n.OutOfRangeException = u, n.OperationException = o
	}, {}],
	21: [function (e, t, n) {
		"use strict";

		function r(e) {
			return e && e.__esModule ? e : {
				default: e
			}
		}

		function i(e, t) {
			if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
		}

		function a(e, t) {
			if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
			return !t || "object" != typeof t && "function" != typeof t ? e : t
		}

		function s(e, t) {
			if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
			e.prototype = Object.create(t && t.prototype, {
				constructor: {
					value: e,
					enumerable: !1,
					writable: !0,
					configurable: !0
				}
			}), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var o = function () {
				function e(e, t) {
					for (var n = 0; n < t.length; n++) {
						var r = t[n];
						r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
					}
				}
				return function (t, n, r) {
					return n && e(t.prototype, n), r && e(t, r), t
				}
			}(),
			u = e("../../scripts/components/events"),
			l = r(u),
			f = e("../../scripts/utils"),
			c = r(f),
			d = e("../../scripts/raise"),
			p = r(d),
			h = e("../../scripts/components/regex"),
			m = r(h),
			g = e("../../scripts/components/array"),
			v = r(g),
			y = e("../../scripts/components/string"),
			b = r(y),
			w = function (e) {
				function t(e, n) {
					var r;
					i(this, t);
					var s = a(this, (t.__proto__ || Object.getPrototypeOf(t)).call(this));
					return s.rules = [], s.searchDisabled = !1, s.init(e, n), r = s, a(s, r)
				}
				return s(t, e), o(t, [{
					key: "init",
					value: function (e, t) {
						t && c.default.extend(this, t), this.options = c.default.extend({}, this.defaults, e)
					}
				}, {
					key: "set",
					value: function (e, t) {
						return t = c.default.extend({
							silent: !1
						}, t || {}), e ? (e.id && !e.key && (e.key = e.id), e.key && (this.rules = c.default.reject(this.rules, function (t) {
							return t.key === e.key
						})), e.fromLiveFilters ? this.setLiveFilter(e) : (this.rules.push(e), t.silent || this.onRulesChange(e), this)) : this
					}
				}, {
					key: "setLiveFilter",
					value: function (e) {
						(0, p.default)(12, "LiveFilter feature not implemented.")
					}
				}, {
					key: "getRuleByKey",
					value: function (e) {
						return c.default.find(this.rules, function (t) {
							return t.key == e
						})
					}
				}, {
					key: "getRulesByType",
					value: function (e) {
						return c.default.where(this.rules, function (t) {
							return t.type == e
						})
					}
				}, {
					key: "removeRuleByKey",
					value: function (e, t) {
						t = c.default.extend({
							silent: !1
						}, t || {});
						var n = this,
							r = n.rules,
							i = c.default.find(r, function (t) {
								return t.key == e
							});
						return i && (n.rules = c.default.reject(r, function (e) {
							return e === i
						}), t.silent || n.onRulesChange()), n
					}
				}, {
					key: "onRulesChange",
					value: function () {}
				}, {
					key: "search",
					value: function (e, t, n) {
						if (!t || !e || this.searchDisabled) return e;
						var r = t instanceof RegExp ? t : m.default.getSearchPattern(b.default.getString(t), n);
						return !!r && (n.searchProperties || (n.searchProperties = this.context.getSearchProperties()), n.searchProperties || (0, p.default)(11, "missing search properties"), v.default.searchByStringProperties({
							pattern: r,
							properties: n.searchProperties,
							collection: e,
							keepSearchDetails: !1
						}))
					}
				}, {
					key: "skim",
					value: function (e) {
						var t = this,
							n = t.rules,
							r = n.length;
						if (!r) return e;
						for (var i = e, a = 0; a < r; a++) {
							var s = t.rules[a];
							s.disabled || (i = t.applyFilter(i, s))
						}
						return i
					}
				}, {
					key: "applyFilter",
					value: function (e, t) {
						switch (t.type) {
							case "search":
								return this.search(e, t.value, t);
							case "fn":
							case "function":
								return c.default.where(e, c.default.partial(t.fn.bind(t.context || this), t))
						}
						return e
					}
				}, {
					key: "reset",
					value: function () {
						for (var e; e = this.rules.shift();) e.onReset && e.onReset.call(this);
						return this
					}
				}], [{
					key: "defaults",
					get: function () {
						return {}
					}
				}]), t
			}(l.default);
		n.default = w
	}, {
		"../../scripts/components/array": 1,
		"../../scripts/components/events": 3,
		"../../scripts/components/regex": 6,
		"../../scripts/components/string": 7,
		"../../scripts/raise": 26,
		"../../scripts/utils": 34
	}],
	22: [function (e, t, n) {
		"use strict";

		function r(e, t) {
			if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
		}

		function i() {
			var e, t, n = arguments.length;
			for (e = 0; e < n; e++)
				if (t = arguments[e], !u.default.isNumber(t)) throw new Error("invalid type")
		}

		function a() {
			var e, t, n = arguments.length;
			for (e = 0; e < n; e++)
				if (t = arguments[e], !u.default.isUnd(t) && !u.default.isNumber(t)) throw new Error("invalid type")
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var s = function () {
				function e(e, t) {
					for (var n = 0; n < t.length; n++) {
						var r = t[n];
						r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
					}
				}
				return function (t, n, r) {
					return n && e(t.prototype, n), r && e(t, r), t
				}
			}(),
			o = e("../../scripts/utils"),
			u = function (e) {
				return e && e.__esModule ? e : {
					default: e
				}
			}(o),
			l = function () {
				function e(t) {
					r(this, e), t = t || {}, a(t.page, t.totalItemsCount, t.resultsPerPage), this.page = t.page || 0, this.resultsPerPage = t.resultsPerPage || 30, this.totalItemsCount = t.totalItemsCount || 1 / 0, this.totalPageCount = 1 / 0, this.firstObjectNumber = void 0, this.lastObjectNumber = void 0, u.default.isUnd(t.totalItemsCount) || this.setTotalItemsCount(t.totalItemsCount, !0), u.default.isFunction(t.onPageChange) && (this.onPageChange = t.onPageChange)
				}
				return s(e, [{
					key: "data",
					value: function () {
						var e = this;
						return {
							page: e._page,
							totalPageCount: e.totalPageCount,
							resultsPerPage: e.resultsPerPage,
							firstObjectNumber: e.firstObjectNumber,
							lastObjectNumber: e.lastObjectNumber,
							totalItemsCount: e.totalItemsCount
						}
					}
				}, {
					key: "prev",
					value: function () {
						var e = this,
							t = e.page - 1;
						return e.validPage(t) && (e.page = t), e
					}
				}, {
					key: "next",
					value: function () {
						var e = this,
							t = e.page + 1;
						return e.validPage(t) && (e.page = t), e
					}
				}, {
					key: "first",
					value: function () {
						return this.page = 1, this
					}
				}, {
					key: "last",
					value: function () {
						return this.page = this.totalPageCount, this
					}
				}, {
					key: "updateItemsNumber",
					value: function () {
						var e = this,
							t = e.totalItemsCount;
						e.firstObjectNumber = e.page * e.resultsPerPage - e.resultsPerPage + 1, e.lastObjectNumber = Math.min(u.default.isNumber(t) ? t : 1 / 0, e.page * e.resultsPerPage)
					}
				}, {
					key: "setTotalItemsCount",
					value: function (e, t) {
						i(e);
						var n = this;
						n.totalItemsCount = e;
						var r = n.getPageCount(e, n.resultsPerPage);
						return n.totalPageCount = r, !t && r < n.page && (n.page = 1), n.updateItemsNumber(), n
					}
				}, {
					key: "validPage",
					value: function (e) {
						var t = this;
						return !(isNaN(e) || e < 1 || e > t.totalPageCount || e === t.page)
					}
				}, {
					key: "onPageChange",
					value: function () {}
				}, {
					key: "getPageCount",
					value: function (e, t) {
						return i(e, t), e === 1 / 0 ? 1 / 0 : e === -1 / 0 ? 0 : e < 1 ? 0 : e > t ? e % t == 0 ? e / t : Math.ceil(e / t) : 1
					}
				}, {
					key: "dispose",
					value: function () {
						delete this.onPageChange
					}
				}, {
					key: "resultsPerPage",
					get: function () {
						return this._resultsPerPage
					},
					set: function (e) {
						e || (e = 0), i(e);
						var t = this,
							n = t.totalItemsCount;
						if (n) {
							var r = t.getPageCount(n, e);
							t.totalPageCount = r, r <= t._page && (t.page = r)
						}
						t._resultsPerPage = e, t.updateItemsNumber()
					}
				}, {
					key: "page",
					get: function () {
						return this._page
					},
					set: function (e) {
						i(e), e != this.page && (this._page = e, this.updateItemsNumber(), this.onPageChange())
					}
				}]), e
			}();
		n.default = l
	}, {
		"../../scripts/utils": 34
	}],
	23: [function (e, t, n) {
		"use strict";

		function r(e, t) {
			if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		}), n.TextSlider = void 0;
		var i = function () {
				function e(e, t) {
					for (var n = 0; n < t.length; n++) {
						var r = t[n];
						r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
					}
				}
				return function (t, n, r) {
					return n && e(t.prototype, n), r && e(t, r), t
				}
			}(),
			a = e("../../scripts/utils"),
			s = (function (e) {
				e && e.__esModule
			}(a), e("../../scripts/exceptions")),
			o = function () {
				function e(t, n) {
					r(this, e), t || (0, s.ArgumentNullException)("text"), this.length = t.length, this.i = 0, this.j = this.length, this.text = t, this.filler = n || " ", this.right = !0
				}
				return i(e, [{
					key: "reset",
					value: function () {
						this.i = 0, this.j = this.length, this.right = !0
					}
				}, {
					key: "next",
					value: function () {
						var e = this,
							t = e.text,
							n = e.filler,
							r = (e.length, e.i),
							i = e.j,
							a = e.right,
							s = t.substr(r, i),
							o = !1;
						return a ? 1 == i ? (i = t.length, r = i, o = !0) : i-- : 0 == r ? (i--, o = !0) : r--, a && t.length != s.length ? s = new Array(t.length - s.length + 1).join(n) + s : s += new Array(t.length - s.length + 1).join(n), o && (a = !a, e.right = a), e.i = r, e.j = i, s
					}
				}]), e
			}();
		n.TextSlider = o
	}, {
		"../../scripts/exceptions": 20,
		"../../scripts/utils": 34
	}],
	24: [function (e, t, n) {
		"use strict";

		function r(e) {
			return e && e.__esModule ? e : {
				default: e
			}
		}

		function i() {
			return new f.VHtmlElement("span", {
				class: "oi",
				"data-glyph": "caret-right"
			})
		}

		function a(e) {
			if (!e) throw "missing menus";
			if (u.default.isPlainObject(e)) return a([e]);
			if (!u.default.isArray(e) || !e.length) throw "missing menus";
			var t = e[0];
			!t.items && t.menu && (e = [{
				items: e
			}]);
			var n = u.default.map(e, function (e) {
				var t = e.items;
				return new f.VHtmlElement("ul", {
					id: e.id,
					class: "ug-menu"
				}, t ? u.default.map(t, function (e) {
					if (e) return s(e)
				}) : null)
			});
			return new f.VWrapperElement(n)
		}

		function s(e) {
			var t, n = e || {},
				r = n.type,
				s = n.href,
				o = n.name,
				l = n.menu,
				c = n.attr,
				d = l ? i() : null,
				p = [],
				h = new f.VTextElement(o || "");
			switch (c && c.css && !c.class && (c.class = c.css, delete c.css), r) {
				case "checkbox":
					var m = u.default.uniqueId("mnck-"),
						g = !!n.checked || void 0;
					t = new f.VWrapperElement([new f.VHtmlElement("input", u.default.extend({}, c, {
						id: m,
						type: "checkbox",
						checked: g
					})), new f.VHtmlElement("label", {
						for: m
					}, h)]);
					break;
				case "radio":
					var v = n.value;
					if (!v) throw new Error("missing 'value' for radio menu item");
					var m = u.default.uniqueId("mnrd-"),
						g = !!n.checked || void 0;
					t = new f.VWrapperElement([new f.VHtmlElement("input", u.default.extend({}, c, {
						id: m,
						type: "radio",
						checked: g,
						value: v
					})), new f.VHtmlElement("label", {
						for: m
					}, h)]);
					break;
				default:
					t = s ? new f.VHtmlElement("a", u.default.extend({
						href: s
					}, c), [h, d]) : new f.VHtmlElement("span", u.default.extend({
						tabindex: "0"
					}, c), [h, d])
			}
			return p.push(t), l && p.push(a(l)), new f.VHtmlElement("li", {
				id: n.id,
				class: l ? "ug-submenu" : void 0
			}, p)
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		}), n.menuItemBuilder = n.menuBuilder = void 0;
		var o = e("../../scripts/utils"),
			u = r(o),
			l = e("../../scripts/dom"),
			f = (r(l), e("../../scripts/exceptions"), e("../../scripts/data/html"));
		n.menuBuilder = a, n.menuItemBuilder = s
	}, {
		"../../scripts/data/html": 12,
		"../../scripts/dom": 19,
		"../../scripts/exceptions": 20,
		"../../scripts/utils": 34
	}],
	25: [function (e, t, n) {
		"use strict";

		function r(e) {
			return e && e.__esModule ? e : {
				default: e
			}
		}

		function i(e) {
			return /input|select|textarea|label|^a$/i.test(e.target.tagName)
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var a = e("../../scripts/dom"),
			s = r(a),
			o = e("../../scripts/utils"),
			u = r(o),
			l = {
				closeMenus: function (e) {
					e && 3 === e.which || u.default.each(["ug-menu", "ug-submenu"], function (t) {
						var n = document.body.getElementsByClassName(t);
						u.default.each(n, function (t) {
							if (!s.default.contains(t, e.target)) {
								var n = t.parentNode;
								s.default.hasClass(n, "open") && (/input|textarea/i.test(e.target.tagName) && s.default.contains(n, e.target) || s.default.removeClass(n, "open"))
							}
						})
					})
				},
				expandMenu: function (e) {
					if (i(e)) return !0;
					var t = e.target;
					if (s.default.hasClass(t, "disabled") || t.hasAttribute("disabled")) return !1;
					var n = t.parentElement;
					return s.default.hasClass(n, "open") ? s.default.removeClass(n, "open") : s.default.addClass(n, "open"), e.preventDefault(), !1
				},
				expandSubMenu: function (e) {
					if (i(e)) return !0;
					var t = s.default.closestWithTag(e.target, "li"),
						n = s.default.siblings(t);
					return u.default.each(n, function (e) {
						s.default.removeClass(e, "open");
						var t = e.getElementsByClassName("open");
						u.default.each(t, function (e) {
							s.default.removeClass(e, "open")
						})
					}), s.default.addClass(t, "open"), !1
				}
			};
		n.default = {
			setup: function () {
				if (this.initialized) return !1;
				this.initialized = !0;
				var e = "click.menus",
					t = document.body;
				s.default.off(t, e), s.default.on(t, e, l.closeMenus), s.default.on(t, e, ".ug-expander", l.expandMenu), s.default.on(t, e, ".ug-submenu", l.expandSubMenu)
			}
		}
	}, {
		"../../scripts/dom": 19,
		"../../scripts/utils": 34
	}],
	26: [function (e, t, n) {
		"use strict";

		function r(e, t) {
			var n = (t || "Error") + ". For further details: https://github.com/RobertoPrevato/KingTable/wiki/Errors#" + e;
			throw "undefined" != typeof console && console.error(n), new Error(n)
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		}), n.default = r
	}, {}],
	27: [function (e, t, n) {
		"use strict";

		function r(e) {
			return e && e.__esModule ? e : {
				default: e
			}
		}

		function i(e, t) {
			if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
		}

		function a(e, t) {
			if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
			return !t || "object" != typeof t && "function" != typeof t ? e : t
		}

		function s(e, t) {
			if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
			e.prototype = Object.create(t && t.prototype, {
				constructor: {
					value: e,
					enumerable: !1,
					writable: !0,
					configurable: !0
				}
			}), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var o = function () {
				function e(e, t) {
					for (var n = 0; n < t.length; n++) {
						var r = t[n];
						r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
					}
				}
				return function (t, n, r) {
					return n && e(t.prototype, n), r && e(t, r), t
				}
			}(),
			u = e("../../scripts/components/events"),
			l = r(u),
			f = e("../../scripts/tables/kingtable"),
			c = r(f),
			d = function (e) {
				function t(e) {
					i(this, t);
					var n = a(this, (t.__proto__ || Object.getPrototypeOf(t)).call(this));
					return n.table = e, n
				}
				return s(t, e), o(t, [{
					key: "getReg",
					value: function () {
						var e = this.table;
						return e ? e.getReg() : c.default.regional.en
					}
				}]), t
			}(l.default);
		n.default = d
	}, {
		"../../scripts/components/events": 3,
		"../../scripts/tables/kingtable": 30
	}],
	28: [function (e, t, n) {
		"use strict";

		function r(e) {
			return e && e.__esModule ? e : {
				default: e
			}
		}

		function i(e, t) {
			if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
		}

		function a(e, t) {
			if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
			return !t || "object" != typeof t && "function" != typeof t ? e : t
		}

		function s(e, t) {
			if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
			e.prototype = Object.create(t && t.prototype, {
				constructor: {
					value: e,
					enumerable: !1,
					writable: !0,
					configurable: !0
				}
			}), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var o = function () {
				function e(e, t) {
					for (var n = 0; n < t.length; n++) {
						var r = t[n];
						r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
					}
				}
				return function (t, n, r) {
					return n && e(t.prototype, n), r && e(t, r), t
				}
			}(),
			u = e("../../scripts/utils"),
			l = r(u),
			f = e("../../scripts/components/string"),
			c = r(f),
			d = e("../../scripts/data/html"),
			p = e("../../scripts/tables/kingtable.builder"),
			h = r(p),
			m = function (e) {
				function t() {
					return i(this, t), a(this, (t.__proto__ || Object.getPrototypeOf(t)).apply(this, arguments))
				}
				return s(t, e), o(t, [{
					key: "getItemAttrObject",
					value: function (e, t) {
						var n = this.options,
							r = n && n.itemDecorator,
							i = {
								class: "kt-item",
								"data-item-ix": e
							};
						if (r) {
							var a = r.call(this, t);
							return l.default.extend(i, a)
						}
						return i
					}
				}, {
					key: "highlight",
					value: function (e, t) {
						if (!e) return "";
						if (!(t instanceof RegExp)) {
							var n = this.table,
								t = n.searchText ? n.filters.getRuleByKey("search").value : null;
							if (!t) return e
						}
						var r, i = c.default.findDiacritics(e),
							a = i.length;
						r = a ? c.default.normalize(e) : e;
						var s = [];
						r.replace(t, function (e) {
							var t = arguments[arguments.length - 2];
							s.push({
								i: t,
								val: a ? c.default.restoreDiacritics(e, i, t) : e
							})
						});
						for (var o, u, l = "", f = 0, p = 0, h = s.length; p < h; p++) {
							o = s[p], u = o.val;
							var m = e.substring(f, o.i);
							l += (0, d.escapeHtml)(m), f = o.i + u.length, l += '<span class="kt-search-highlight">' + (0, d.escapeHtml)(u) + "</span>"
						}
						return f < e.length && (l += (0, d.escapeHtml)(e.substr(f))), l
					}
				}]), t
			}(h.default);
		n.default = m
	}, {
		"../../scripts/components/string": 7,
		"../../scripts/data/html": 12,
		"../../scripts/tables/kingtable.builder": 27,
		"../../scripts/utils": 34
	}],
	29: [function (e, t, n) {
		"use strict";

		function r(e) {
			return e && e.__esModule ? e : {
				default: e
			}
		}

		function i(e, t) {
			if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
		}

		function a(e, t) {
			if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
			return !t || "object" != typeof t && "function" != typeof t ? e : t
		}

		function s(e, t) {
			if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
			e.prototype = Object.create(t && t.prototype, {
				constructor: {
					value: e,
					enumerable: !1,
					writable: !0,
					configurable: !0
				}
			}), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var o = function () {
				function e(e, t) {
					for (var n = 0; n < t.length; n++) {
						var r = t[n];
						r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
					}
				}
				return function (t, n, r) {
					return n && e(t.prototype, n), r && e(t, r), t
				}
			}(),
			u = e("../../scripts/data/html"),
			l = e("../../scripts/tables/kingtable.html.base.builder"),
			f = r(l),
			c = e("../../scripts/raise"),
			d = r(c),
			p = e("../../scripts/utils"),
			h = r(p),
			m = e("../../scripts/dom"),
			g = r(m),
			v = function (e) {
				function t(e) {
					i(this, t);
					var n = a(this, (t.__proto__ || Object.getPrototypeOf(t)).call(this, e));
					return n.options = h.default.extend({}, e ? e.options : null), n.setListeners(), n
				}
				return s(t, e), o(t, [{
					key: "setListeners",
					value: function () {
						var e = this,
							t = e.table;
						if (!t || !t.element) return e;
						e.listenTo(t, {
							"fetching:data": function () {
								e.loadingHandler()
							},
							"fetched:data": function () {
								e.unsetLoadingHandler()
							},
							"fetch:fail": function () {
								e.unsetLoadingHandler().display(e.errorView())
							},
							"no-results": function () {
								e.unsetLoadingHandler().display(e.emptyView())
							}
						})
					}
				}, {
					key: "getGeneratedFields",
					value: function () {
						var e = this.options,
							t = this.getReg(),
							n = e.detailRoute,
							r = [],
							i = t.goToDetails;
						if (n) {
							/\/$/.test(n) || (n = e.detailRoute = n + "/");
							var a = this.table.getIdProperty();
							r.push({
								name: "details-link",
								html: function (e) {
									return "<a class='kt-details-link' href='" + (n + e[a]) + "'>" + i + "</a>"
								}
							})
						}
						return r
					}
				}, {
					key: "getFields",
					value: function () {
						var e = this.table,
							t = h.default.clone(e.columns),
							n = this.options,
							r = n.itemCount,
							i = this.getGeneratedFields(),
							a = r ? {
								name: "ε_row",
								displayName: "#"
							} : null;
						a && t.unshift(a);
						var s = n.fields;
						return s ? (h.default.isFunction(s) && (s = s.call(this, t)), t = i.concat(s.concat(t))) : t = i.concat(t), t
					}
				}, {
					key: "build",
					value: function () {
						var e = this,
							t = e.table,
							n = t.getData({
								format: !0,
								hide: !1
							});
						if (!n || !n.length) return e.display(e.emptyView());
						var r = e.getFields(),
							i = e.buildCaption(),
							a = e.buildView(r, n),
							s = e.buildRoot(i, a);
						e.display(s)
					}
				}, {
					key: "buildRoot",
					value: function (e, t) {
						var n = this.table,
							r = {
								class: "king-table-region"
							};
						return n.id && (r.id = n.id), new u.VHtmlElement("div", r, [e, t])
					}
				}, {
					key: "buildView",
					value: function (e, t) {
						this.table;
						return new u.VHtmlElement("table", {
							class: "king-table"
						}, [this.buildHead(e), this.buildBody(e, t)])
					}
				}, {
					key: "buildHead",
					value: function (e) {
						var t = (this.table, new u.VHtmlElement("tr", {}, h.default.map(h.default.values(e), function (e) {
							if (!e.hidden && !e.secret) return new u.VHtmlElement("th", {
								class: e.css
							}, new u.VTextElement(e.displayName))
						})));
						return new u.VHtmlElement("thead", {
							class: "king-table-head"
						}, t)
					}
				}, {
					key: "buildBody",
					value: function (e, t) {
						var n = this,
							r = this.table,
							i = r.builder,
							a = r.options.formattedSuffix,
							s = r.searchText ? r.filters.getRuleByKey("search").value : null,
							o = r.options.autoHighlightSearchProperties,
							l = -1,
							f = h.default.map(t, function (t) {
								l += 1, t.__ix__ = l;
								for (var r, f, c = [], p = 0, m = e.length; p < m; p++)
									if (f = e[p], r = f.name, !f.hidden && !f.secret) {
										var g, v = r + a,
											y = h.default.has(t, v) ? t[v] : t[r];
										if (f.html) {
											h.default.isFunction(f.html) || (0, d.default)(24, "Invalid 'html' option for property, it must be a function.");
											var b = f.html.call(i, t, y);
											g = new u.VHtmlFragment(b || "")
										} else g = null === y || void 0 === y || "" === y ? new u.VTextElement("") : s && o && h.default.isString(y) ? new u.VHtmlFragment(i.highlight(y, s)) : new u.VTextElement(y);
										c.push(new u.VHtmlElement("td", f ? {
											class: f.css || f.name
										} : {}, g))
									}
								return new u.VHtmlElement("tr", n.getItemAttrObject(l, t), c)
							});
						return new u.VHtmlElement("tbody", {
							class: "king-table-body"
						}, f)
					}
				}, {
					key: "buildCaption",
					value: function () {
						var e = this.table,
							n = e.options.caption,
							r = t.options.paginationInfo ? this.buildPaginationInfo() : null;
						return n || r ? new u.VHtmlElement("div", {
							class: "king-table-caption"
						}, [n ? new u.VHtmlElement("span", {}, new u.VTextElement(n)) : null, r && n ? new u.VHtmlElement("br") : null, r]) : null
					}
				}, {
					key: "buildPaginationInfo",
					value: function () {
						var e = this.table,
							t = e.pagination,
							n = this.getReg(),
							r = t.page,
							i = t.totalPageCount,
							a = (t.resultsPerPage, t.firstObjectNumber),
							s = t.lastObjectNumber,
							o = t.totalItemsCount,
							l = e.getFormattedAnchorTime(),
							f = h.default.isNumber,
							c = "";
						return f(r) && (c += n.page + " " + r, f(i) && i > 0 && (c += " " + n.of + " " + i), f(a) && f(s) && s > 0 && (c += " - " + n.results + " " + a + " - " + s, f(o) && (c += " " + n.of + " - " + o))), l && e.options.showAnchorTimestamp && (c += " - " + n.anchorTime + " " + l), new u.VHtmlElement("span", {
							class: "pagination-info"
						}, new u.VTextElement(c))
					}
				}, {
					key: "emptyView",
					value: function (e) {
						var t = this.getReg(),
							n = new u.VHtmlElement("div", {
								class: "king-table-empty"
							}, new u.VHtmlElement("span", 0, new u.VTextElement(t.noData)));
						return e ? n : this.singleLine(this.table, n)
					}
				}, {
					key: "errorView",
					value: function (e) {
						return e || (e = this.getReg().errorFetchingData), this.singleLine(this.table, new u.VHtmlFragment('<div class="king-table-error">\n      <span class="message">\n        <span>' + e + '</span>\n        <span class="oi" data-glyph="warning" aria-hidden="true"></span>\n      </span>\n    </div>'))
					}
				}, {
					key: "loadingView",
					value: function () {
						var e = (this.table, this.getReg()),
							t = this.buildCaption();
						return t.children.push(new u.VHtmlElement("div", {
							class: "loading-info"
						}, [new u.VHtmlElement("span", {
							class: "loading-text"
						}, new u.VTextElement(e.loading)), new u.VHtmlElement("span", {
							class: "mini-loader"
						})])), this.buildRoot([t])
					}
				}, {
					key: "display",
					value: function (e) {
						var t = this.table;
						h.default.isString(e) || (e = e.toString());
						var n = t.element;
						if (n) {
							for (n.classList.add("king-table"), t.emit("empty:element", n); n.hasChildNodes();) n.removeChild(n.lastChild);
							n.innerHTML = e
						}
					}
				}, {
					key: "singleLine",
					value: function (e) {
						var t = (this.table, this.buildCaption());
						return t.children.push(new u.VHtmlElement("br"), new u.VHtmlElement("div", {
							class: "loading-info"
						}, h.default.isString(e) ? new u.VTextElement(e) : e)), this.buildRoot([t])
					}
				}, {
					key: "loadingHandler",
					value: function () {
						var e = this,
							n = e.table;
						e.unsetLoadingHandler();
						var r = n.hasData() ? t.options.loadInfoDelay : 0;
						e.showLoadingTimeout = setTimeout(function () {
							if (!n.loading) return e.unsetLoadingHandler();
							e.display(e.loadingView())
						}, r)
					}
				}, {
					key: "unsetLoadingHandler",
					value: function () {
						return clearTimeout(this.showLoadingTimeout), this.showLoadingTimeout = null, this
					}
				}, {
					key: "dispose",
					value: function () {
						var e = this.table,
							t = e.element;
						t && g.default.empty(t), this.stopListening(this.table), this.table = null, delete this.options
					}
				}], [{
					key: "options",
					get: function () {
						return {
							handleLoadingInfo: !0,
							loadInfoDelay: 500,
							paginationInfo: !0
						}
					}
				}]), t
			}(f.default);
		n.default = v
	}, {
		"../../scripts/data/html": 12,
		"../../scripts/dom": 19,
		"../../scripts/raise": 26,
		"../../scripts/tables/kingtable.html.base.builder": 28,
		"../../scripts/utils": 34
	}],
	30: [function (e, t, n) {
		"use strict";

		function r(e) {
			return e && e.__esModule ? e : {
				default: e
			}
		}

		function i(e, t) {
			if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
		}

		function a(e, t) {
			if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
			return !t || "object" != typeof t && "function" != typeof t ? e : t
		}

		function s(e, t) {
			if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
			e.prototype = Object.create(t && t.prototype, {
				constructor: {
					value: e,
					enumerable: !1,
					writable: !0,
					configurable: !0
				}
			}), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var o = function () {
				function e(e, t) {
					for (var n = 0; n < t.length; n++) {
						var r = t[n];
						r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
					}
				}
				return function (t, n, r) {
					return n && e(t.prototype, n), r && e(t, r), t
				}
			}(),
			u = "function" == typeof Symbol && "symbol" == typeof Symbol.iterator ? function (e) {
				return typeof e
			} : function (e) {
				return e && "function" == typeof Symbol && e.constructor === Symbol && e !== Symbol.prototype ? "symbol" : typeof e
			},
			l = e("../../scripts/tables/kingtable.text.builder"),
			f = r(l),
			c = e("../../scripts/tables/kingtable.html.builder"),
			d = r(c),
			p = e("../../scripts/tables/kingtable.html.base.builder"),
			h = r(p),
			m = e("../../scripts/tables/kingtable.rhtml.builder"),
			g = r(m),
			v = e("../../scripts/tables/kingtable.regional"),
			y = r(v),
			b = e("../../scripts/components/events"),
			w = r(b),
			k = e("../../scripts/data/object-analyzer"),
			x = r(k),
			E = e("../../scripts/data/sanitizer"),
			S = r(E),
			P = e("../../scripts/filters/filters-manager"),
			T = r(P),
			_ = e("../../scripts/filters/paginator"),
			C = r(_),
			O = e("../../scripts/data/ajax"),
			V = r(O),
			j = e("../../scripts/raise"),
			M = r(j),
			F = e("../../scripts/utils"),
			D = r(F),
			H = e("../../scripts/components/string"),
			I = r(H),
			N = e("../../scripts/components/regex"),
			A = r(N),
			L = e("../../scripts/components/number"),
			R = r(L),
			B = e("../../scripts/components/date"),
			z = r(B),
			U = e("../../scripts/components/array"),
			q = r(U),
			W = e("../../scripts/data/csv"),
			$ = r(W),
			K = e("../../scripts/data/json"),
			Y = r(K),
			Z = e("../../scripts/data/xml"),
			J = r(Z),
			X = e("../../scripts/data/file"),
			G = r(X),
			Q = e("../../scripts/data/lru"),
			ee = r(Q),
			te = e("../../scripts/data/memstore"),
			ne = r(te),
			re = e("../../scripts/exceptions"),
			ie = {
				lang: "en",
				caption: null,
				itemCount: !0,
				columnDefault: {
					name: "",
					type: "text",
					sortable: !0,
					allowSearch: !0,
					hidden: !1
				},
				httpMethod: "GET",
				allowSearch: !0,
				minSearchChars: 3,
				page: 1,
				resultsPerPage: 30,
				formattedSuffix: "_(formatted)",
				fixed: void 0,
				search: "",
				sortByFormatter: q.default.humanSortBy,
				searchMode: "FullString",
				exportFormats: [{
					name: "Csv",
					format: "csv",
					type: "text/csv",
					cs: !0
				}, {
					name: "Json",
					format: "json",
					type: "application/json",
					cs: !0
				}, {
					name: "Xml",
					format: "xml",
					type: "text/xml",
					cs: !0
				}],
				prettyXml: !0,
				csvOptions: {},
				exportHiddenProperties: !1,
				builder: "rhtml",
				storeTableData: !0,
				lruCacheSize: 10,
				lruCacheMaxAge: 9e5,
				showAnchorTimestamp: !0,
				collectionName: "data",
				searchSortingRules: !0,
				idProperty: null,
				autoHighlightSearchProperties: !0,
				emptyValue: ""
			},
			ae = {
				text: f.default,
				html: d.default,
				rhtml: g.default
			};
		if ("undefined" == ("undefined" == typeof Promise ? "undefined" : u(Promise))) {
			var se = !1;
			if ("undefined" != ("undefined" == typeof ES6Promise ? "undefined" : u(ES6Promise)) && ES6Promise.polyfill) try {
				ES6Promise.polyfill(), se = "undefined" != ("undefined" == typeof Promise ? "undefined" : u(Promise))
			} catch (e) {}
			se || (0, M.default)(1, "Missing implementation of Promise (missing dependency)")
		}
		var oe = function (e) {
			function t(e, n) {
				i(this, t);
				var r = a(this, (t.__proto__ || Object.getPrototypeOf(t)).call(this));
				e = e || {};
				var s = r;
				n && D.default.extend(s, n), D.default.each(s.baseProperties(), function (t) {
					D.default.has(e, t) && (s[t] = e[t], delete e[t])
				});
				var o = e.sortBy;
				return D.default.isString(o) && (s.sortCriteria = q.default.getSortCriteria(o)), e = s.options = D.default.extend({}, t.defaults, e), s.loading = !1, s.init(e, n), r
			}
			return s(t, e), o(t, [{
				key: "baseProperties",
				value: function () {
					return ["id", "onInit", "element", "context", "fixed", "prepareData", "getExtraFilters", "getTableData", "afterRender", "beforeRender", "numberFilterFormatter", "dateFilterFormatter"]
				}
			}, {
				key: "init",
				value: function (e) {
					var t = this;
					t.cache = {}, t.disposables = [], t.analyst = new x.default, t.sanitizer = new S.default, t.filters = new T.default({}, {
						context: t
					});
					var n = e.data;
					return n && (D.default.isArray(n) || (0, M.default)(3, "Data is not an array"), n = Y.default.clone(n), t.setFixedData(n)), t.loadSettings(), t.setPagination(), t.fixed || (t.filters.searchDisabled = !0), t.setBuilder(e.builder).onInit(), t
				}
			}, {
				key: "getReg",
				value: function () {
					var e = this.options.lang;
					e || (0, M.default)(15, "Missing language option (cannot be null or falsy)");
					var n = t.regional[e];
					return n || (0, M.default)(15, "Missing regional for language " + e), n
				}
			}, {
				key: "setBuilder",
				value: function (e) {
					e || (0, M.default)(8, "name cannot be null or empty");
					var n = this,
						r = (n.options, t.builders);
					n.builder && n.disposeOf(n.builder);
					var i = r[e];
					i || (0, M.default)(10, "Missing handler for builder: " + e);
					var a = new i(n);
					return n.builder = a, n.disposables.push(a), n
				}
			}, {
				key: "getSubSet",
				value: function (e) {
					var t = this.pagination,
						n = (t.page - 1) * t.resultsPerPage,
						r = t.resultsPerPage + n;
					return e.slice(n, r)
				}
			}, {
				key: "setPagination",
				value: function () {
					var e = this,
						t = e.data,
						n = e.options,
						r = n.page,
						i = n.resultsPerPage,
						a = n.totalItemsCount || (t ? t.length : 0);
					e.pagination && e.disposeOf(e.pagination);
					var s = e.pagination = new C.default({
						page: r,
						resultsPerPage: i,
						totalItemsCount: a,
						onPageChange: function () {
							e.render()
						}
					});
					return e.disposables.push(s), e
				}
			}, {
				key: "updatePagination",
				value: function (e) {
					var t = this;
					if (t.pagination || t.setPagination(), !D.default.isNumber(e)) throw "invalid type";
					return t.pagination.setTotalItemsCount(e), D.default.ifcall(t.onResultsCountChange, t), t.trigger("change:pagination"), t
				}
			}, {
				key: "onInit",
				value: function () {}
			}, {
				key: "hasData",
				value: function () {
					var e = this.data;
					return !(!e || !e.length)
				}
			}, {
				key: "getItemStructure",
				value: function () {
					return this.analyst.describe(this.data, {
						lazy: !0
					})
				}
			}, {
				key: "initColumns",
				value: function () {
					var e = "columnsInitialized",
						n = this;
					if (n[e] || !n.hasData()) return n;
					n[e] = !0;
					var r, i = [],
						a = n.getItemStructure(),
						s = [];
					for (r in a) a[r] = {
						name: r,
						type: a[r]
					}, s.push(r);
					var o = n.options.columns;
					if (o) {
						var u, l = 0;
						for (r in o) {
							var f = o[r];
							D.default.isPlainObject(o) ? u = r : D.default.isArray(o) ? (D.default.isString(f) && (0, M.default)(16, "invalid columns option " + f), (u = f.name) || (0, M.default)(17, "missing name in column option")) : (0, M.default)(16, "invalid columns option"), -1 == D.default.indexOf(s, u) && (0, M.default)(18, 'A column is defined with name "' + u + '", but this property was not found among object properties. Items properties are: ' + s.join(", ")), D.default.isString(f) && (f = o[r] = {
								displayName: f
							}), D.default.isString(f.name) && (f.displayName = f.name, delete f.name), a[u] = D.default.extend(a[u], f);
							var c = f.position;
							D.default.isNumber(c) ? a[u].position = c : a[u].position = l, l++
						}
					}
					for (r in a) {
						var d = {
								name: r
							},
							p = a[r],
							h = p.type;
						h || (p.type = h = "string");
						var f = D.default.extend({}, n.options.columnDefault, d, p);
						f.cid = D.default.uniqueId("col"), h = D.default.lower(h);
						var m = t.Schemas.DefaultByType;
						if (D.default.has(m, h)) {
							var g = m[h];
							D.default.isFunction(g) && (g = g.call(n, p, a)), D.default.extend(d, g)
						}
						if (m = t.Schemas.DefaultByName, D.default.has(m, r) && D.default.extend(d, m[r]), D.default.extend(f, d), o) {
							var v = D.default.isArray(o) ? D.default.find(o, function (e) {
								return e.name == r
							}) : o[r];
							v && D.default.extend(f, v)
						}
						e = "css", D.default.isString(f[e]) || (f[e] = I.default.kebabCase(f.name)), D.default.isString(f.displayName) || (f.displayName = f.name), i.push(f)
					}
					var y = "position";
					if (o) {
						var l = 0;
						for (var r in o) {
							var f = D.default.find(i, function (e) {
								return e.name == r
							});
							f && !D.default.has(f, y) && (f[y] = l), l++
						}
					}
					var b = n.getCachedColumnsData(),
						w = "hidden";
					return b && D.default.each(b, function (e) {
						var t = D.default.find(i, function (t) {
							return t.name == e.name
						});
						t && (t[y] = e[y], t[w] = e[w])
					}), n.columns = i, n.fixed && n.searchText && n.setSearchFilter(n.searchText, !0), n
				}
			}, {
				key: "storePreference",
				value: function (e, t) {
					var n = this.getFiltersStore();
					if (!n) return !1;
					var r = this.getMemoryKey(e);
					n.setItem(r, t)
				}
			}, {
				key: "getPreference",
				value: function (e) {
					var t = this.getFiltersStore();
					if (t) {
						var n = this.getMemoryKey(e);
						return t.getItem(n)
					}
				}
			}, {
				key: "getFiltersStore",
				value: function () {
					return localStorage
				}
			}, {
				key: "getDataStore",
				value: function () {
					return sessionStorage
				}
			}, {
				key: "loadSettings",
				value: function () {
					return this.getTableData && this.restoreTableData(), this.restoreFilters()
				}
			}, {
				key: "restoreFilters",
				value: function () {
					var e = this,
						t = e.options,
						n = e.getFiltersStore();
					if (!n) return e;
					var r = e.getMemoryKey("filters"),
						i = n.getItem(r);
					if (!i) return e;
					try {
						var a = Y.default.parse(i),
							s = "page size sortBy search timestamp".split(" ");
						e.trigger("restore:filters", a), D.default.each(s, function (n) {
							"search" == n ? e.validateForSeach(a[n]) && e.setSearchFilter(a[n], !0) : "sortBy" == n ? e.sortCriteria = a[n] : "size" == n ? t.resultsPerPage = a[n] : t[n] = a[n]
						});
						var o = D.default.minus(a, s);
						D.default.isEmpty(o) || e.restoreExtraFilters(a)
					} catch (e) {
						n.removeItem(r)
					}
					return e
				}
			}, {
				key: "getFilters",
				value: function () {
					var e = this,
						t = e.pagination;
					return D.default.isUnd(e.anchorTime) && (e.anchorTime = new Date), D.default.extend({}, e.getExtraFilters(), {
						page: t.page,
						size: t.resultsPerPage,
						sortBy: e.sortCriteria || null,
						search: e.searchText || null,
						timestamp: e.anchorTime || null
					})
				}
			}, {
				key: "getFiltersSetCache",
				value: function () {
					var e = this,
						t = e.getFilters(),
						n = e.getFiltersStore();
					if (n) {
						var r = e.getMemoryKey("filters");
						e.trigger("store:filters", t), n.setItem(r, Y.default.compose(t))
					}
					return t
				}
			}, {
				key: "getExtraFilters",
				value: function () {}
			}, {
				key: "restoreExtraFilters",
				value: function (e) {}
			}, {
				key: "beforeRender",
				value: function () {}
			}, {
				key: "afterRender",
				value: function () {}
			}, {
				key: "onFetchStart",
				value: function () {}
			}, {
				key: "onFetchDone",
				value: function () {}
			}, {
				key: "onFetchFail",
				value: function () {}
			}, {
				key: "onFetchEnd",
				value: function () {}
			}, {
				key: "onSearchEmpty",
				value: function () {}
			}, {
				key: "onSearchStart",
				value: function (e) {}
			}, {
				key: "prepareData",
				value: function (e) {
					return this
				}
			}, {
				key: "formatValues",
				value: function (e) {
					var t = this,
						n = t.options;
					e || (e = t.data);
					var r, i, a = t.options.formattedSuffix,
						s = D.default.where(t.columns, function (e) {
							return D.default.isFunction(e.format)
						});
					return D.default.each(e, function (e) {
						D.default.each(s, function (t) {
							r = t.name + a, i = e[t.name], e[r] = D.default.isUnd(i) || null === i || "" === i ? n.emptyValue : t.format(i, e) || n.emptyValue
						})
					}), t
				}
			}, {
				key: "setColumnsOrder",
				value: function () {
					var e = D.default.stringArgs(arguments),
						t = e.length;
					if (!t) return !1;
					for (var n = this.columns, r = [], i = 0; i < t; i++) {
						var a = e[i],
							s = D.default.find(n, function (e) {
								return e.name == a
							});
						s || (0, M.default)(19, 'missing column with name "' + a + '"'), s.position = i, r.push(s)
					}
					var o = D.default.where(n, function (e) {
						return -1 == r.indexOf(e)
					});
					return D.default.each(o, function (e) {
						t++, e.position = t
					}), q.default.sortBy(n, "position"), this.storeColumnsData().render(), this
				}
			}, {
				key: "toggleColumns",
				value: function (e) {
					var t = this.columns;
					return D.default.each(e, function (e) {
						if (D.default.isArray(e)) var n = e[0],
							r = e[1];
						else var n = e.name,
							r = e.visible;
						var i = D.default.find(t, function (e) {
							return e.name == n
						});
						i || (0, M.default)(19, 'missing column with name "' + n + '"'), i.hidden = !r || !!i.secret
					}), this.storeColumnsData().render(), this
				}
			}, {
				key: "hideColumns",
				value: function () {
					return this.columnsVisibility(D.default.stringArgs(arguments), !1)
				}
			}, {
				key: "showColumns",
				value: function () {
					return this.columnsVisibility(D.default.stringArgs(arguments), !0)
				}
			}, {
				key: "columnsVisibility",
				value: function (e, t) {
					var n = this;
					return 1 == e.length && "*" == e[0] ? D.default.each(this.columns, function (e) {
						e.hidden = !t || !!e.secret
					}) : D.default.each(e, function (e) {
						n.colAttr(e, "hidden", !t)
					}), this.storeColumnsData().render(), this
				}
			}, {
				key: "colAttr",
				value: function (e, t, n) {
					e || (0, re.ArgumentNullException)("name");
					var r = this.columns;
					r || (0, M.default)(20, "missing columns information (properties not initialized)");
					var i = D.default.find(r, function (t) {
						return t.name == e
					});
					return i || (0, M.default)(19, 'missing column with name "' + e + '"'), i[t] = n, i
				}
			}, {
				key: "storeColumnsData",
				value: function () {
					var e = this.getDataStore(),
						t = this.getMemoryKey("columns:data"),
						n = this.options,
						r = this.columns;
					return e && n.storeTableData && e.setItem(t, Y.default.compose(r)), this
				}
			}, {
				key: "getCachedColumnsData",
				value: function () {
					var e = this.getDataStore(),
						t = this.getMemoryKey("columns:data"),
						n = this.options;
					if (e && n.storeTableData) {
						var r = e.getItem(t);
						if (r) try {
							return Y.default.parse(r)
						} catch (n) {
							e.removeItem(t)
						}
					}
					return null
				}
			}, {
				key: "sortColumns",
				value: function () {
					if (arguments.length) return this.setColumnsOrder.apply(this, arguments);
					var e = D.default.isNumber,
						t = this.columns;
					t.sort(function (t, n) {
						var r = "position";
						return e(t[r]) && !e(n[r]) ? -1 : !e(t[r]) && e(n[r]) ? 1 : t[r] > n[r] ? 1 : t[r] < n[r] ? -1 : (r = "displayName", I.default.compare(t[r], n[r], 1))
					});
					for (var n = 0, r = t.length; n < r; n++) t[n].position = n;
					return this
				}
			}, {
				key: "setTools",
				value: function () {
					return this
				}
			}, {
				key: "handleTableData",
				value: function (e) {
					this.tableData = e
				}
			}, {
				key: "refresh",
				value: function () {
					return delete this.anchorTime, this.render({
						clearDataCache: !0
					})
				}
			}, {
				key: "hardRefresh",
				value: function () {
					return this.trigger("hard:refresh").clearTableData(), this.render({
						clearDataCache: !0
					})
				}
			}, {
				key: "render",
				value: function (e) {
					var t = this;
					return new Promise(function (n, r) {
						function i() {
							t.beforeRender(), t.initColumns().sortColumns().setTools().build(), t.afterRender(), n()
						}

						function a() {
							if (t.fixed && t.hasData()) t.getFiltersSetCache(), i();
							else {
								var n = t.lastFetchTimestamp = (new Date).getTime();
								t.getList(e, n).then(function (e) {
									if (!e || !e.length && !t.columnsInitialized) return t.emit("no-results");
									i()
								}, function () {
									t.emit("get-list:failed"), r("get-list:failed")
								})
							}
						}
						if (t.getTableData && !t.cache.tableDataFetched) {
							var s = t.getTableData();
							D.default.quacksLikePromise(s) || (0, M.default)(13, "getTableData must return a Promise or Promise-like object."), s.then(function (e) {
								t.options.storeTableData && (t.cache.tableDataFetched = !0, t.storeTableData(e)), t.handleTableData(e), a()
							}, function () {
								t.emit("get-table-data:failed"), r()
							})
						} else a()
					})
				}
			}, {
				key: "clearTableData",
				value: function () {
					var e = this,
						t = e.getDataStore(),
						n = e.getMemoryKey("table:data"),
						r = e.options;
					return t && r.storeTableData && t.removeItem(n), e.cache.tableDataFetched = !1, delete e.tableData, e
				}
			}, {
				key: "storeTableData",
				value: function (e) {
					var t = this.getDataStore(),
						n = this.getMemoryKey("table:data"),
						r = this.options;
					return t && r.storeTableData && t.setItem(n, Y.default.compose(e)), this
				}
			}, {
				key: "restoreTableData",
				value: function () {
					var e = this,
						t = e.getDataStore(),
						n = e.getMemoryKey("table:data"),
						r = e.options;
					if (t && r.storeTableData) {
						var i = t.getItem(n);
						if (i) {
							try {
								i = Y.default.parse(i)
							} catch (e) {
								t.removeItem(n)
							}
							e.handleTableData(i), e.cache.tableDataFetched = !0
						}
					}
					return e
				}
			}, {
				key: "build",
				value: function () {
					var e = this,
						t = e.builder;
					t && t.build()
				}
			}, {
				key: "getMemoryKey",
				value: function (e) {
					var t = location.pathname + location.hash + ".kt",
						n = this.id;
					return n && (t = n + ":" + t), e ? t + ":::" + e : t
				}
			}, {
				key: "setFixedData",
				value: function (e) {
					var t = this;
					return e = t.normalizeCollection(e), t.fixed = !0, t.filters.searchDisabled = !1, t.prepareData(e), t.data = e, t.initColumns(), t.formatValues(e), t.updatePagination(e.length), e
				}
			}, {
				key: "getList",
				value: function (e, t) {
					e = e || {};
					var n = this,
						r = n.mixinFetchData();
					return new Promise(function (i, a) {
						n.emit("fetch:start").onFetchStart(), n.getFetchPromiseWithCache(r, e).then(function (e) {
							if (!(t < n.lastFetchTimestamp))
								if (e || (0, M.default)(14, "`getFetchPromise` did not return a value when resolving"), n.emit("fetch:done").onFetchDone(e), D.default.isArray(e)) e = n.setFixedData(e), i(n.getSubSet(e));
								else {
									var r = e.items || e.subset;
									D.default.isArray(r) || (0, M.default)(6, "The returned object is not a catalog"), D.default.isNumber(e.total) || (0, M.default)(7, "Missing total items count in response object."), r = n.normalizeCollection(r), n.prepareData(r), n.data = r, n.initColumns(), n.formatValues(r), n.updatePagination(e.total), i(r)
								}
						}, function () {
							t < n.lastFetchTimestamp || (n.emit("fetch:fail").onFetchFail(), a())
						}).then(function () {
							n.emit("fetch:end").onFetchEnd()
						})
					})
				}
			}, {
				key: "search",
				value: function (e) {
					D.default.isUnd(e) && (e = "");
					var t = this;
					t.validateForSeach(e) ? (e ? (t.onSearchStart(e), t.setSearchFilter(e)) : t.unsetSearch(), t.pagination.page = 1) : t.unsetSearch(), t.render()
				}
			}, {
				key: "isSearchActive",
				value: function () {
					return !!this.filters.getRuleByKey("search")
				}
			}, {
				key: "unsetSearch",
				value: function () {
					var e = this;
					return e.isSearchActive() ? (e.filters.removeRuleByKey("search"), e.searchText = null, e.hasData() && e.updatePagination(e.data.length), e.trigger("search-empty").onSearchEmpty(), e) : e
				}
			}, {
				key: "setSearchFilter",
				value: function (e, t) {
					var n = this;
					n.searchText = e, t || n.getFiltersSetCache();
					var r = n.getSearchProperties();
					return n.filters.set({
						type: "search",
						key: "search",
						value: A.default.getSearchPattern(I.default.getString(e), {
							searchMode: n.options.searchMode
						}),
						searchProperties: !(!r || !r.length) && r
					}), n.trigger("search-active"), n
				}
			}, {
				key: "getSearchProperties",
				value: function () {
					var e = this,
						t = e.options;
					if (t.searchProperties) return t.searchProperties;
					if (!e.data || !e.columnsInitialized) return !1;
					var n = D.default.where(e.columns, function (e) {
							return e.allowSearch && !e.secret
						}),
						r = t.formattedSuffix;
					return D.default.flatten(D.default.map(n, function (e) {
						return D.default.isFunction(e.format) ? "number" == e.type ? [e.name + r, e.name] : e.name + r : e.name
					}))
				}
			}, {
				key: "validateForSeach",
				value: function (e) {
					if (!e) return !1;
					var t = this.options.minSearchChars;
					return !(e.match(/^[\s]+$/g) || D.default.isNumber(t) && e.length < t)
				}
			}, {
				key: "getFormattedAnchorTime",
				value: function () {
					var e = this.anchorTime;
					return e instanceof Date ? z.default.isToday(e) ? z.default.format(e, "HH:mm:ss") : z.default.formatWithTime(e) : ""
				}
			}, {
				key: "getFormattedFetchTime",
				value: function () {
					var e = this.dataFetchTime;
					return e instanceof Date ? z.default.isToday(e) ? z.default.format(e, "HH:mm:ss") : z.default.formatWithTime(e) : ""
				}
			}, {
				key: "getFetchPromiseWithCache",
				value: function (e, t) {
					t || (t = {});
					var n = this,
						r = n.options,
						i = r.lruCacheSize,
						a = n.getDataStore(),
						s = !(!i || !a);
					if (s) {
						var o = Y.default.parse(Y.default.compose(e)),
							u = n.getMemoryKey("catalogs"),
							l = ee.default.get(u, function (e) {
								return D.default.equal(o, e.filters)
							}, a, !0);
						if (l) {
							if (!t.clearDataCache) return n.anchorTime = new Date(l.data.anchorTime), n.dataFetchTime = new Date(l.ts), new Promise(function (e, t) {
								setTimeout(function () {
									e(l.data.data)
								}, 0)
							});
							ee.default.remove(u, void 0, a)
						}
					}
					return new Promise(function (t, i) {
						n.loading = !0, n.emit("fetching:data"), n.getFetchPromise(e).then(function (i) {
							s && ee.default.set(u, {
								data: i,
								filters: e,
								anchorTime: n.anchorTime.getTime()
							}, r.lruCacheSize, r.lruCacheMaxAge, a), n.dataFetchTime = new Date, n.loading = !1, n.emit("fetched:data"), t(i)
						}, function () {
							n.loading = !1, i()
						})
					})
				}
			}, {
				key: "getFetchPromise",
				value: function (e) {
					var t = this.options,
						n = t.url;
					n || (0, M.default)(5, "Missing url option, to fetch data");
					var r = t.httpMethod;
					return e = this.formatFetchData(e), V.default.shot({
						type: r,
						url: n,
						data: e
					})
				}
			}, {
				key: "numberFilterFormatter",
				value: function (e, t) {
					return t
				}
			}, {
				key: "dateFilterFormatter",
				value: function (e, t) {
					return z.default.toIso8601(t)
				}
			}, {
				key: "formatFetchData",
				value: function (e) {
					var t = this.options,
						n = t.sortByFormatter;
					e.sortBy && D.default.isFunction(n) && (e.sortBy = n(e.sortBy));
					var r;
					for (r in e) {
						var i = e[r];
						i instanceof Date && (e[r] = this.dateFilterFormatter(r, i)), i instanceof Number && (e[r] = this.numberFilterFormatter(r, i))
					}
					return e
				}
			}, {
				key: "mixinFetchData",
				value: function () {
					var e = this.options.fetchData;
					return D.default.isFunction(e) && (e = e.call(this)), D.default.extend(this.getFiltersSetCache(), e || {})
				}
			}, {
				key: "normalizeCollection",
				value: function (e) {
					var t = e.length;
					if (!t) return e;
					var n = e[0];
					if (D.default.isArray(n)) {
						var r, i, a = [],
							s = n.length;
						for (r = 1; r < t; r++) {
							var o = {};
							for (i = 0; i < s; i++) o[n[i]] = e[r][i];
							a.push(o)
						}
						return a
					}
					return e
				}
			}, {
				key: "getData",
				value: function (e) {
					var t = D.default.extend({
							optimize: !1,
							itemCount: !0,
							hide: !0
						}, e),
						n = this,
						r = n.options.itemCount && t.itemCount,
						i = n.getItemsToDisplay(),
						a = D.default.clone(n.columns);
					return t.hide && (D.default.each(D.default.where(n.columns, function (e) {
						return e.hidden || e.secret
					}), function (e) {
						D.default.each(i, function (t) {
							delete t[e.name]
						})
					}), a = D.default.where(n.columns, function (e) {
						return !e.hidden && !e.secret
					})), r && n.setItemsNumber(i), t.optimize ? (r && a.unshift({
						name: "ε_row",
						displayName: "#"
					}), n.optimizeCollection(i, a, t)) : i
				}
			}, {
				key: "getItemsToDisplay",
				value: function () {
					var e = this,
						t = e.options,
						n = e.data;
					if (!n || !n.length) return [];
					if (n = D.default.clone(n), !e.fixed) return n;
					var r = n.length;
					if (n = e.filters.skim(n), n.length != r && e.updatePagination(n.length), !e.searchText || !t.searchSortingRules) {
						var i = e.sortCriteria;
						D.default.isEmpty(i) || q.default.sortBy(n, i)
					}
					return e.getSubSet(n)
				}
			}, {
				key: "sortBy",
				value: function () {
					var e = q.default.getSortCriteria(arguments);
					if (!e || !e.length) return this.unsetSortBy();
					var t = this;
					return t.sortCriteria = e, t.hasData() ? (q.default.sortBy(t.data, e), t.render()) : t.getFiltersSetCache(), t
				}
			}, {
				key: "progressSortBy",
				value: function (e) {
					e || (0, re.ArgumentNullException)("name");
					var t = this,
						n = t.columns;
					n || (0, M.default)(20, "Missing columns information"), D.default.find(n, function (t) {
						return t.name == e
					}) || (0, M.default)(19, "Column '${name}' is not found among columns.");
					var r = t.sortCriteria || [],
						i = D.default.find(r, function (t) {
							return t[0] == e
						});
					if (i) {
						-1 === i[1] ? r = D.default.reject(r, function (t) {
							return t[0] == e
						}) : i[1] = -1
					} else r.push([e, 1]);
					t.sortBy(r)
				}
			}, {
				key: "progressSortBySingle",
				value: function (e) {
					e || (0, re.ArgumentNullException)("name");
					var t = this,
						n = t.columns;
					n || (0, M.default)(20, "Missing columns information"), D.default.find(n, function (t) {
						return t.name == e
					}) || (0, M.default)(19, "Column '${name}' is not found among columns.");
					var r = t.sortCriteria || [],
						i = D.default.find(r, function (t) {
							return t[0] == e
						});
					if (i) {
						r = -1 === i[1] ? [e, 1] : [e, -1]
					} else r = [e, 1];
					t.sortBy([r])
				}
			}, {
				key: "unsetSortBy",
				value: function () {
					return this.sortCriteria = null, this.render(), this
				}
			}, {
				key: "setItemsNumber",
				value: function (e) {
					var t = this,
						n = t.pagination,
						r = (n.page - 1) * n.resultsPerPage;
					e || (e = t.data);
					for (var i = e.length, a = 0; a < i; a++) e[a].ε_row = (a + 1 + r).toString();
					return e
				}
			}, {
				key: "optimizeCollection",
				value: function (e, t, n) {
					t || (t = this.columns), n || (n = {
						format: !0
					});
					for (var r, i = [D.default.map(t, function (e) {
							return e.displayName
						})], a = n.format, s = this.options.formattedSuffix, o = 0, u = e.length; o < u; o++) {
						for (var l = [], f = 0, c = t.length; f < c; f++) {
							var d = t[f].name,
								p = d + s;
							r = e[o], a && D.default.has(r, p) ? l.push(r[p]) : l.push(r[d] || "")
						}
						i.push(l)
					}
					return i
				}
			}, {
				key: "getItemValue",
				value: function (e, t) {
					e || (0, re.ArgumentNullException)("item");
					var n = this.options,
						r = n.formattedSuffix,
						i = t + r;
					return D.default.has(e, i) ? e[i] : e[t]
				}
			}, {
				key: "getIdProperty",
				value: function () {
					var e = this.options;
					if (D.default.isString(e.idProperty)) return e.idProperty;
					var t = this.columns;
					t && t.length || (0, M.default)(4, "id property cannot be determined: columns are not initialized.");
					for (var n = 0, r = t.length; n < r; n++) {
						var i = t[n].name;
						if (/^_?id$|^_?guid$/i.test(i)) return i
					}(0, M.default)(4, "id property cannot be determined, please specify it using 'idProperty' option.")
				}
			}, {
				key: "getExportFileName",
				value: function (e) {
					return this.options.collectionName + "." + e
				}
			}, {
				key: "getColumnsForExport",
				value: function () {
					var e = this.columns;
					return this.options.exportHiddenProperties ? e : D.default.reject(e, function (e) {
						return e.hidden || e.secret
					})
				}
			}, {
				key: "exportTo",
				value: function (e) {
					e || (0, re.ArgumentException)("format");
					var t = this,
						n = t.options,
						r = t.getExportFileName(e),
						i = D.default.find(t.options.exportFormats, function (t) {
							return t.format === e
						});
					t.getColumnsForExport();
					i && i.type || (0, M.default)(30, "Missing format information");
					var a = t.getData({
							itemCount: !1
						}),
						s = "";
					if (i.handler) s = i.handler.call(t, a);
					else switch (e) {
						case "csv":
							var o = t.optimizeCollection(a);
							s = $.default.serialize(o, n.csvOptions);
							break;
						case "json":
							s = Y.default.compose(a, 2, 2);
							break;
						case "xml":
							s = t.dataToXml(a);
							break;
						default:
							throw "export format " + e + "not implemented"
					}
					s && G.default.exportfile(r, s, i.type)
				}
			}, {
				key: "dataToXml",
				value: function (e) {
					for (var t = this, n = t.getColumnsForExport(), r = t.options, i = document, a = new XMLSerializer, s = i.createElement(r.collectionName || "collection"), o = 0, u = e.length; o < u; o++) {
						for (var l = i.createElement(r.entityName || "item"), f = 0, c = n.length; f < c; f++) {
							var d = n[f],
								p = d.name,
								h = e[o][p];
							if (r.entityUseProperties) l.setAttribute(p, h);
							else {
								var m = i.createElement(p);
								m.innerText = h, l.appendChild(m)
							}
						}
						s.appendChild(l)
					}
					var g = a.serializeToString(s);
					return r.prettyXml ? J.default.pretty(g) : J.default.normal(g)
				}
			}, {
				key: "disposeOf",
				value: function (e) {
					e.dispose(), D.default.removeItem(this.disposables, e)
				}
			}, {
				key: "dispose",
				value: function () {
					delete this.context, delete this.search, delete this.filters.context, D.default.each(this.disposables, function (e) {
						e.dispose && e.dispose(), D.default.isFunction(e) && e()
					}), this.disposables = [];
					var e = this.options;
					D.default.ifcall(e.onDispose, this)
				}
			}], [{
				key: "regional",
				get: function () {
					return y.default
				}
			}, {
				key: "version",
				get: function () {
					return "2.0.0"
				}
			}, {
				key: "Utils",
				get: function () {
					return D.default
				}
			}, {
				key: "StringUtils",
				get: function () {
					return I.default
				}
			}, {
				key: "NumberUtils",
				get: function () {
					return R.default
				}
			}, {
				key: "ArrayUtils",
				get: function () {
					return q.default
				}
			}, {
				key: "DateUtils",
				get: function () {
					return z.default
				}
			}, {
				key: "json",
				get: function () {
					return Y.default
				}
			}, {
				key: "Paginator",
				get: function () {
					return C.default
				}
			}, {
				key: "PlainTextBuilder",
				get: function () {
					return f.default
				}
			}, {
				key: "HtmlBuilder",
				get: function () {
					return d.default
				}
			}, {
				key: "RichHtmlBuilder",
				get: function () {
					return g.default
				}
			}, {
				key: "BaseHtmlBuilder",
				get: function () {
					return h.default
				}
			}, {
				key: "builders",
				get: function () {
					return ae
				}
			}, {
				key: "stores",
				get: function () {
					return {
						memory: ne.default
					}
				}
			}]), t
		}(w.default);
		oe.defaults = ie, oe.Schemas = {
			DefaultByType: {
				number: function (e, t) {
					return {
						format: function (e) {
							return R.default.format(e)
						}
					}
				},
				date: function (e, t) {
					return {
						format: function (e) {
							var t = oe.DateUtils.hasTime(e),
								n = oe.DateUtils.defaults.format[t ? "long" : "short"];
							return oe.DateUtils.format(e, n)
						}
					}
				}
			},
			DefaultByName: {
				id: {
					name: "id",
					type: "id",
					hidden: !0,
					secret: !0
				},
				guid: {
					name: "guid",
					type: "guid",
					hidden: !0,
					secret: !0
				}
			}
		}, oe.Ajax = V.default, "undefined" !== ("undefined" == typeof window ? "undefined" : u(window)) && (window.KingTable = oe), n.default = oe
	}, {
		"../../scripts/components/array": 1,
		"../../scripts/components/date": 2,
		"../../scripts/components/events": 3,
		"../../scripts/components/number": 4,
		"../../scripts/components/regex": 6,
		"../../scripts/components/string": 7,
		"../../scripts/data/ajax": 9,
		"../../scripts/data/csv": 10,
		"../../scripts/data/file": 11,
		"../../scripts/data/json": 13,
		"../../scripts/data/lru": 14,
		"../../scripts/data/memstore": 15,
		"../../scripts/data/object-analyzer": 16,
		"../../scripts/data/sanitizer": 17,
		"../../scripts/data/xml": 18,
		"../../scripts/exceptions": 20,
		"../../scripts/filters/filters-manager": 21,
		"../../scripts/filters/paginator": 22,
		"../../scripts/raise": 26,
		"../../scripts/tables/kingtable.html.base.builder": 28,
		"../../scripts/tables/kingtable.html.builder": 29,
		"../../scripts/tables/kingtable.regional": 31,
		"../../scripts/tables/kingtable.rhtml.builder": 32,
		"../../scripts/tables/kingtable.text.builder": 33,
		"../../scripts/utils": 34
	}],
	31: [function (e, t, n) {
		"use strict";
		Object.defineProperty(n, "__esModule", {
			value: !0
		}), n.default = {
			en: {
				goToDetails: "Go to details",
				sortOptions: "Sort options",
				searchSortingRules: "When searching, sort by relevance",
				advancedFilters: "Advanced filters",
				sortModes: {
					simple: "Simple (single property)",
					complex: "Complex (multiple properties)"
				},
				viewsType: {
					table: "Table",
					gallery: "Gallery"
				},
				exportFormats: {
					csv: "Csv",
					json: "Json",
					xml: "Xml"
				},
				columns: "Columns",
				export: "Export",
				view: "View",
				views: "Views",
				loading: "Loading",
				noData: "No data to display",
				page: "Page",
				resultsPerPage: "Results per page",
				results: "Results",
				of: "of",
				firstPage: "First page",
				lastPage: "Last page",
				prevPage: "Previous page",
				nextPage: "Next page",
				refresh: "Refresh",
				fetchTime: "Data fetched at:",
				anchorTime: "Data at:",
				sortAscendingBy: "Sort by {{name}} ascending",
				sortDescendingBy: "Sort by {{name}} descending",
				errorFetchingData: "An error occurred while fetching data."
			}
		}
	}, {}],
	32: [function (e, t, n) {
		"use strict";

		function r(e) {
			return e && e.__esModule ? e : {
				default: e
			}
		}

		function i(e, t) {
			if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
		}

		function a(e, t) {
			if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
			return !t || "object" != typeof t && "function" != typeof t ? e : t
		}

		function s(e, t) {
			if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
			e.prototype = Object.create(t && t.prototype, {
				constructor: {
					value: e,
					enumerable: !1,
					writable: !0,
					configurable: !0
				}
			}), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
		}

		function o(e) {
			return {
				class: e
			}
		}

		function u(e) {
			return e || (0, w.default)(34, "Invalid extra view configuration."), e.name || (0, w.default)(35, "Missing name in extra view configuration."), e.getItemTemplate && (x.default.extend(e, {
				resolver: {
					getItemTemplate: e.getItemTemplate,
					buildView: function (t, n, r) {
						var i = this.getItemTemplate();
						i || (0, w.default)(31, "Invalid getItemTemplate function in extra view.");
						var a = x.default.map(r, function (e) {
							var n = i.replace(/\{\{(.+?)\}\}/g, function (n, r) {
								return e.hasOwnProperty(r) || (0, w.default)(32, "Missing property " + r + ", for template"), t.getItemValue(e, r)
							});
							return new c.VHtmlFragment(n)
						});
						return new c.VHtmlElement("div", {
							class: ("king-table-body " + e.name).toLowerCase()
						}, a)
					}
				}
			}), delete e.getItemTemplate), e
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var l = function e(t, n, r) {
				null === t && (t = Function.prototype);
				var i = Object.getOwnPropertyDescriptor(t, n);
				if (void 0 === i) {
					var a = Object.getPrototypeOf(t);
					return null === a ? void 0 : e(a, n, r)
				}
				if ("value" in i) return i.value;
				var s = i.get;
				if (void 0 !== s) return s.call(r)
			},
			f = function () {
				function e(e, t) {
					for (var n = 0; n < t.length; n++) {
						var r = t[n];
						r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
					}
				}
				return function (t, n, r) {
					return n && e(t.prototype, n), r && e(t, r), t
				}
			}(),
			c = e("../../scripts/data/html"),
			d = e("../../scripts/menus/kingtable.menu.html"),
			p = (r(d), e("../../scripts/tables/kingtable.html.builder")),
			h = r(p),
			m = e("../../scripts/tables/kingtable.html.base.builder"),
			g = r(m),
			v = e("../../scripts/menus/kingtable.menu"),
			y = r(v),
			b = e("../../scripts/raise"),
			w = r(b),
			k = e("../../scripts/utils"),
			x = r(k),
			E = e("../../scripts/dom"),
			S = r(E),
			P = e("../../scripts/data/file"),
			T = r(P),
			_ = function (e) {
				function t() {
					return i(this, t), a(this, (t.__proto__ || Object.getPrototypeOf(t)).apply(this, arguments))
				}
				return s(t, e), f(t, [{
					key: "buildView",
					value: function (e, t, n) {
						return new c.VHtmlElement("div", o("king-table-gallery"), [this.buildBody(e, t, n), new c.VHtmlElement("br", o("break"))])
					}
				}, {
					key: "buildBody",
					value: function (e, t, n) {
						var r = this,
							i = e.options.formattedSuffix,
							a = e.searchText ? e.filters.getRuleByKey("search").value : null,
							s = e.options.autoHighlightSearchProperties,
							o = -1,
							u = x.default.map(n, function (e) {
								o += 1, e.__ix__ = o;
								for (var n, u, l = [], f = 0, d = t.length; f < d; f++)
									if (u = t[f], n = u.name, !u.hidden && !u.secret) {
										var p, h = n + i,
											m = x.default.has(e, h) ? e[h] : e[n];
										if (u.html) {
											x.default.isFunction(u.html) || (0, w.default)(24, "Invalid 'html' option for property");
											var g = u.html.call(r, e, m);
											p = new c.VHtmlFragment(g || "")
										} else p = null === m || void 0 === m || "" === m ? new c.VTextElement("") : a && s && x.default.isString(m) ? new c.VHtmlFragment(r.highlight(m, a)) : new c.VTextElement(m);
										l.push(new c.VHtmlElement("ε_row" == u.name ? "strong" : "span", u ? {
											class: u.css || u.name,
											title: u.displayName
										} : {}, p))
									}
								return new c.VHtmlElement("li", r.getItemAttrObject(o, e), l)
							});
						return new c.VHtmlElement("ul", {
							class: "king-table-body"
						}, u)
					}
				}]), t
			}(g.default),
			C = {
				Simple: "simple",
				Complex: "complex"
			},
			O = function (e) {
				function t(e) {
					i(this, t);
					var n = a(this, (t.__proto__ || Object.getPrototypeOf(t)).call(this, e));
					n.options = x.default.extend({}, t.defaults, e.options, e.options.rhtml || e.options.html), n.setSeachHandler();
					var r = n.options,
						s = r.extraViews;
					return s && (r.views = r.views.concat(x.default.map(s, function (e) {
						return u(e)
					}))), n.loadSettings(), y.default.initialized || y.default.setup(), n.filtersViewOpen = r.filtersView && r.filtersViewExpandable && r.filtersViewOpen, n
				}
				return s(t, e), f(t, [{
					key: "setListeners",
					value: function () {
						var e = this;
						l(t.prototype.__proto__ || Object.getPrototypeOf(t.prototype), "setListeners", this).call(this);
						var n = this,
							r = n.table;
						if (!r || !r.element) return n;
						n.listenTo(r, {
							"change:pagination": function () {
								if (!e.rootElement) return !0;
								n.updatePagination()
							},
							"get-list:failed": function () {
								if (!e.rootElement) return !0;
								n.updatePagination()
							}
						})
					}
				}, {
					key: "loadSettings",
					value: function () {
						var e = this,
							t = e.options,
							n = e.table;
						if (!n.getFiltersStore()) return e;
						var r = n.getPreference("sort-mode");
						r && (t.sortMode = r);
						var i = n.getPreference("view-type");
						return i && (t.view = i), e
					}
				}, {
					key: "setSeachHandler",
					value: function () {
						function e(e) {
							var t = this.table;
							return t.validateForSeach(e) ? t.search(e) : t.isSearchActive() && (t.unsetSearch(), t.render()), t.getFiltersSetCache(), !0
						}
						var t = this.options.searchDelay;
						return this.search = x.default.isNumber(t) && t > 0 ? x.default.debounce(e, t, this) : e, this
					}
				}, {
					key: "getViewResolver",
					value: function () {
						var e = this.options,
							t = e.view,
							n = e.views;
						x.default.isString(t) || (0, w.default)(21, "Missing view configuration for Rich HTML builder");
						var r = x.default.find(n, function (e) {
							return e.name == t
						});
						r || (0, w.default)(22, "Missing view resolver for view: " + t);
						var i = r.resolver;
						return !0 === i ? this : (i || (0, w.default)(33, "Missing resolver in view configuration '" + t + "'"), x.default.isPlainObject(i) || (i = new i), x.default.quacks(i, ["buildView"]) || (0, w.default)(23, "Invalid resolver for view: " + t), i)
					}
				}, {
					key: "buildCaption",
					value: function () {
						var e = this.table,
							t = e.options.caption;
						return t ? new c.VHtmlElement("div", {
							class: "king-table-caption"
						}, new c.VHtmlElement("span", {}, new c.VTextElement(t))) : null
					}
				}, {
					key: "build",
					value: function () {
						var e = this;
						return e.table.element ? e.ensureLayout().update() : e
					}
				}, {
					key: "ensureLayout",
					value: function () {
						var e = this;
						if (e.rootElement) return e;
						var t = e.table,
							n = e.options,
							r = t.element,
							i = e.buildView(null, null, new c.VHtmlFragment(" ")),
							a = e.buildCaption(),
							s = e.buildRoot(a, i);
						return t.emit("empty:element", r), S.default.empty(r), S.default.addClass(r, "king-table"), r.innerHTML = s.toString(), e.rootElement = S.default.findFirstByClass(r, "king-table-region"), e.bindEvents(), x.default.ifcall(n.onLayoutRender, e, [r]), n.filtersView && x.default.ifcall(n.onFiltersRender, e, [S.default.findFirstByClass(r, "kt-filters")]), e
					}
				}, {
					key: "update",
					value: function () {
						this.updatePagination().updateView()
					}
				}, {
					key: "updatePagination",
					value: function () {
						var e = this.table,
							t = e.pagination,
							n = this.rootElement;
						n || (0, w.default)(26, "missing root element");
						var r = this.getReg(),
							t = (e.options, e.pagination),
							i = t.page,
							a = t.totalPageCount,
							s = t.resultsPerPage,
							o = t.firstObjectNumber,
							u = t.lastObjectNumber,
							l = t.totalItemsCount,
							f = e.getFormattedAnchorTime(),
							c = x.default.isNumber,
							d = S.default.findFirstByClass,
							p = S.default.addClass,
							h = S.default.removeClass;
						d(n, "pagination-bar-page-number").value = i, d(n, "pagination-bar-results-select").value = s;
						var m = "pagination-button",
							g = "pagination-button-disabled";
						x.default.each(["pagination-bar-first-page", "pagination-bar-prev-page"], function (e) {
							var t = d(n, e);
							i > 1 ? (p(t, m), h(t, g)) : (p(t, g), h(t, m))
						}), x.default.each(["pagination-bar-last-page", "pagination-bar-next-page"], function (e) {
							var t = d(n, e);
							i < a ? (p(t, m), h(t, g)) : (p(t, g), h(t, m))
						});
						var v = "";
						c(o) && c(u) && u > 0 && (v += r.results + " " + o + " - " + u, c(l) && (v += " " + r.of + " - " + l));
						var y = "";
						f && e.options.showAnchorTimestamp && (y = r.anchorTime + " " + f);
						var b, k = {
							"results-info": v,
							"anchor-timestamp-info": y,
							"total-page-count": r.of + " " + a
						};
						for (b in k) {
							var E = d(n, b);
							E && (E.innerHTML = k[b])
						}
						var P = e.searchText || "",
							T = d(n, "search-field");
						return T && T.value != P && 0 == S.default.isFocused(T) && (T.value = P), this
					}
				}, {
					key: "updateView",
					value: function () {
						var e = this,
							t = e.options,
							n = e.table,
							r = n.pagination,
							i = e.rootElement;
						i || (0, w.default)(26, "missing root element"), x.default.each({
							"kt-search-active": n.searchText,
							"kt-search-sorting": n.options.searchSortingRules
						}, function (e, t) {
							S.default.modClass(i, t, e)
						});
						var r = n.getData({
								format: !0,
								hide: !1
							}),
							a = S.default.findFirstByClass(i, "king-table-view");
						if (!r || !r.length) return a.innerHTML = e.emptyView().toString(), e;
						var s = e.getFields();
						if (e._must_build_tools) {
							document.getElementById(e.toolsRegionId).innerHTML = e.buildToolsInner(!0), delete e._must_build_tools
						}
						e.currentItems = r;
						var o = e.buildView(s, r);
						return a.innerHTML = o.children[0].toString(), x.default.ifcall(t.onViewUpdate, e, [a]), e
					}
				}, {
					key: "display",
					value: function (e) {
						var t = (this.table, this.options);
						x.default.isString(e) || (e = e.toString()), this.ensureLayout();
						var n = this.rootElement,
							r = S.default.findFirstByClass(n, "king-table-view");
						r.innerHTML = e, x.default.ifcall(t.onViewUpdate, this, [r])
					}
				}, {
					key: "buildRoot",
					value: function (e, t) {
						var n = this.table,
							r = {
								class: "king-table-region"
							};
						return n.id && (r.id = n.id), new c.VHtmlElement("div", r, [e, this.buildPaginationBar(), this.buildFiltersView(), t])
					}
				}, {
					key: "buildPaginationBar",
					value: function () {
						var e = this.table,
							t = this.getReg(),
							n = this.options,
							r = e.pagination,
							i = (r.page, r.totalPageCount, r.resultsPerPage, r.firstObjectNumber),
							a = r.lastObjectNumber,
							s = r.totalItemsCount,
							o = n.filtersView,
							u = o && n.filtersViewExpandable,
							l = u && n.filtersViewOpen,
							f = e.getFormattedAnchorTime(),
							d = x.default.isNumber,
							p = "";
						d(i) && d(a) && a > 0 && (p += t.results + " " + i + " - " + a, d(s) && (p += " " + t.of + " - " + s));
						var h;
						f && e.options.showAnchorTimestamp && (h = t.anchorTime + " " + f);
						var m = t.advancedFilters,
							g = n.allowSearch ? new c.VHtmlElement("span", {
								class: "pagination-bar-filters"
							}, new c.VHtmlElement("input", {
								type: "text",
								class: "search-field",
								value: e.searchText || ""
							})) : null,
							v = "span",
							y = new c.VHtmlElement(v, {
								class: "separator"
							});
						return new c.VHtmlElement("div", {
							class: "pagination-bar"
						}, [this.buildTools(), new c.VHtmlElement(v, {
							class: "pagination-bar-buttons"
						}, [new c.VHtmlElement(v, {
							tabindex: "0",
							class: "pagination-button pagination-bar-first-page oi",
							"data-glyph": "media-step-backward",
							title: t.firstPage
						}), new c.VHtmlElement(v, {
							tabindex: "0",
							class: "pagination-button pagination-bar-prev-page oi",
							"data-glyph": "caret-left",
							title: t.prevPage
						}), y, new c.VHtmlElement(v, {
							class: "valigned"
						}, new c.VTextElement(t.page)), new c.VHtmlElement("input", {
							type: "text",
							name: "page-number",
							class: "must-integer pagination-bar-page-number",
							value: r.page
						}), new c.VHtmlElement("span", {
							class: "valigned total-page-count",
							value: r.page
						}, new c.VTextElement(t.of + " " + r.totalPageCount)), y, new c.VHtmlElement(v, {
							tabindex: "0",
							class: "pagination-button pagination-bar-refresh oi",
							"data-glyph": "reload",
							title: t.refresh
						}), y, new c.VHtmlElement(v, {
							tabindex: "0",
							class: "pagination-button pagination-bar-next-page oi",
							"data-glyph": "caret-right",
							title: t.nextPage
						}), new c.VHtmlElement(v, {
							tabindex: "0",
							class: "pagination-button pagination-bar-last-page oi",
							"data-glyph": "media-step-forward",
							title: t.lastPage
						}), y, new c.VHtmlElement(v, {
							class: "valigned"
						}, new c.VTextElement(t.resultsPerPage)), new c.VHtmlElement("select", {
							name: "pageresults",
							class: "pagination-bar-results-select valigned"
						}, x.default.map(n.resultsPerPageSelect, function (e) {
							var t = new c.VHtmlElement("option", {
								value: e
							}, new c.VTextElement(e.toString()));
							return e === n.resultsPerPage && (t.attributes.selected = !0), t
						})), y, p ? new c.VHtmlElement(v, {
							class: "valigned results-info"
						}, new c.VTextElement(p)) : null, p ? y : null, h ? new c.VHtmlElement(v, {
							class: "valigned anchor-timestamp-info"
						}, new c.VTextElement(h)) : null, g ? y : null, g, u ? y : null, u ? new c.VHtmlElement("button", {
							class: "btn valigned camo-btn kt-advanced-filters" + (l ? " kt-open" : "")
						}, new c.VTextElement(m)) : null])])
					}
				}, {
					key: "buildHead",
					value: function (e) {
						var t = this.table,
							n = t.builder,
							r = t.sortCriteria,
							i = n.getReg(),
							a = new c.VHtmlElement("tr", {}, x.default.map(x.default.values(e), function (e) {
								if (!e.hidden && !e.secret) {
									var t, n = !1,
										a = [e.css];
									if (e.sortable) {
										a.push("sortable");
										var s = x.default.find(r, function (t) {
											return t[0] === e.name
										});
										s && (n = !0, t = s[1])
									}
									var o = e.displayName;
									return new c.VHtmlElement("th", {
										class: a.join(" "),
										"data-prop": e.name
									}, new c.VHtmlElement("div", {}, [new c.VHtmlElement("span", {}, new c.VTextElement(o)), n ? new c.VHtmlElement("span", {
										class: "oi kt-sort-glyph",
										"data-glyph": 1 == t ? "sort-ascending" : "sort-descending",
										"aria-hidden": !0,
										title: x.default.format(1 == t ? i.sortAscendingBy : i.sortDescendingBy, {
											name: o
										})
									}) : null]))
								}
							}));
						return new c.VHtmlElement("thead", {
							class: "king-table-head"
						}, a)
					}
				}, {
					key: "buildView",
					value: function (e, t, n) {
						var r, i = this.table;
						if (n) r = n;
						else if (t && t.length) {
							var r, a = this.getViewResolver();
							a === this ? r = new c.VHtmlElement("table", {
								class: "king-table"
							}, [this.buildHead(e), this.buildBody(e, t)]) : (a.table = this.table, a.options = i.options, r = a.buildView(i, e, t), delete a.table, delete a.options)
						} else r = new c.VHtmlElement("div", {
							class: "king-table-view"
						}, this.emptyView());
						return new c.VHtmlElement("div", {
							class: "king-table-view"
						}, r)
					}
				}, {
					key: "getTemplate",
					value: function (e, t) {
						if (x.default.isFunction(e)) return e.call(this);
						x.default.isString(e) || (0, w.default)(38, "Cannot obtain HTML from given parameter " + t + ", must be a function or a string.");
						var n = document.getElementById(e);
						if (null != n) {
							if (/script/i.test(n.tagName)) return n.innerText;
							(0, w.default)(38, "Cannot obtain HTML from parameter " + t + ". Element is not <script>.")
						}
						return e
					}
				}, {
					key: "buildFiltersView",
					value: function () {
						var e = this,
							t = e.options,
							n = t.filtersView;
						if (n) {
							var r = t.filtersViewOpen,
								i = t.filtersViewExpandable,
								a = e.getTemplate(n, "filtersView"),
								s = ["kt-filters"];
							return !r && i || s.push("kt-open"), i && s.push("kt-expandable"), new c.VHtmlElement("div", {
								class: s.join(" ")
							}, [new c.VHtmlFragment(a)])
						}
					}
				}, {
					key: "buildTools",
					value: function () {
						var e = this.table,
							t = e.columnsInitialized;
						t || (this._must_build_tools = !0);
						var n = this.toolsRegionId = x.default.uniqueId("tools-region");
						return new c.VHtmlElement("div", {
							id: n,
							class: "tools-region"
						}, this.buildToolsInner(t))
					}
				}, {
					key: "buildToolsInner",
					value: function (e) {
						return new c.VWrapperElement([new c.VHtmlElement("span", {
							class: "oi ug-expander",
							tabindex: "0",
							"data-glyph": "cog"
						}), e ? this.buildMenu() : null])
					}
				}, {
					key: "buildMenu",
					value: function () {
						var e = this,
							t = e.options,
							n = t.tools,
							r = [e.getColumnsMenuSchema(), e.getViewsMenuSchema(), t.allowSortModes ? e.getSortModeSchema() : null, e.getExportMenuSchema()];
						return n && (x.default.isFunction(n) && (n = n.call(this)), n && (x.default.isArray(n) || (0, w.default)(40, "Tools is not an array or a function returning an array."), r = r.concat(n))), t.prepTools && (x.default.isFunction(t.prepTools) || (0, w.default)(41, "prepTools option must be a function."), t.prepTools.call(this, r)), (0, d.menuBuilder)({
							items: r
						})
					}
				}, {
					key: "getSortModeSchema",
					value: function () {
						var e = this.getReg(),
							t = this.options,
							n = t.sortMode,
							r = x.default.map(C, function (t, r) {
								return {
									name: e.sortModes[r],
									checked: n == r,
									type: "radio",
									value: r,
									attr: {
										name: "kt-sort-mode",
										class: "sort-mode-radio"
									}
								}
							});
						return {
							name: e.sortOptions,
							menu: {
								items: r
							}
						}
					}
				}, {
					key: "getColumnsMenuSchema",
					value: function () {
						if (!this.table.columns || !this.table.columns.length) throw "Columns not initialized.";
						var e = x.default.where(this.table.columns, function (e) {
							return !e.secret
						});
						return {
							name: this.getReg().columns,
							menu: {
								items: x.default.map(e, function (e) {
									return {
										name: e.displayName,
										checked: !e.hidden,
										type: "checkbox",
										attr: {
											name: e.name,
											class: "visibility-check"
										}
									}
								})
							}
						}
					}
				}, {
					key: "getViewsMenuSchema",
					value: function () {
						var e = this.getReg(),
							t = this.options,
							n = t.views,
							r = t.view,
							i = x.default.map(n, function (t) {
								var n = t.name;
								return {
									name: e.viewsType[n] || n,
									checked: r == n,
									type: "radio",
									value: n,
									attr: {
										name: "kt-view-type",
										class: "view-type-radio"
									}
								}
							});
						return {
							name: e.view,
							menu: {
								items: i
							}
						}
					}
				}, {
					key: "getExportMenuSchema",
					value: function () {
						var e = this.table,
							t = e.options.exportFormats;
						if (!t || !t.length) return null;
						if (T.default.supportsCsExport() || (t = x.default.reject(t, function (e) {
								return e.cs || e.clientSide
							})), !t || !t.length) return null;
						var n = this.getReg(),
							r = x.default.map(t, function (e) {
								return {
									name: n.exportFormats[e.format],
									attr: {
										css: "export-btn",
										"data-format": e.format
									}
								}
							});
						return {
							name: n.export,
							menu: {
								items: r
							}
						}
					}
				}, {
					key: "goToPrev",
					value: function () {
						this.table.pagination.prev()
					}
				}, {
					key: "goToNext",
					value: function () {
						this.table.pagination.next()
					}
				}, {
					key: "goToFirst",
					value: function () {
						this.table.pagination.first()
					}
				}, {
					key: "goToLast",
					value: function () {
						this.table.pagination.last()
					}
				}, {
					key: "refresh",
					value: function () {
						this.table.refresh()
					}
				}, {
					key: "changePage",
					value: function (e) {
						var t = e.target.value;
						/^\d+$/.test(t) && this.table.pagination.validPage(parseInt(t)) ? (this.table.pagination.page = parseInt(t), this.table.render()) : e.target.value = this.table.pagination.page
					}
				}, {
					key: "changeResultsNumber",
					value: function (e) {
						var t = e.target.value;
						this.table.pagination.resultsPerPage = parseInt(t), this.table.render()
					}
				}, {
					key: "getItemByEv",
					value: function (e, t) {
						if (e) return this.getItemByEl(e.target, t)
					}
				}, {
					key: "getItemByEl",
					value: function (e, t) {
						if (e) {
							var n = S.default.closestWithClass(e, "kt-item");
							if (!n) {
								if (t) return;
								(0, w.default)(36, "Cannot retrieve an item by event data. Make sure that HTML elements generated for table items have 'kt-item' class.")
							}
							var r = n.dataset.itemIx;
							return x.default.isUnd(r) && (0, w.default)(37, "Cannot retrieve an item by element data. Make sure that HTML elements generated for table items have 'data-ix' attribute."), this.currentItems[r]
						}
					}
				}, {
					key: "onItemClick",
					value: function (e) {
						var t = this.getItemByEl(e.target),
							n = this.options,
							r = n.purist;
						n.onItemClick.call(this, t, r ? void 0 : e)
					}
				}, {
					key: "toggleAdvancedFilters",
					value: function () {
						var e = "filtersViewOpen",
							t = "kt-open",
							n = S.default.findByClass(this.rootElement, "kt-filters")[0],
							r = S.default.hasClass(n, t);
						this[e] = !r, S.default.modClass(n, t, this[e])
					}
				}, {
					key: "clearFilters",
					value: function () {}
				}, {
					key: "sort",
					value: function (e) {
						var t = e.target,
							n = this.options;
						if (this.table.searchText && n.searchSortingRules) return !0;
						/th/i.test(t.tagName) || (t = S.default.closestWithTag(t, "th"));
						var r = t.dataset.prop,
							i = this.table;
						if (r && x.default.any(this.table.columns, function (e) {
								return e.name == r
							})) switch (n.sortMode) {
							case C.Simple:
								i.progressSortBySingle(r);
								break;
							case C.Complex:
								i.progressSortBy(r);
								break;
							default:
								(0, w.default)(28, "Invalid sort mode options. Value must be either 'simple' or 'complex'.")
						}
					}
				}, {
					key: "onSearchKeyUp",
					value: function (e) {
						var t = e.target.value;
						this.search(t)
					}
				}, {
					key: "onSearchChange",
					value: function (e) {
						var t = e.target.value;
						this.search(t)
					}
				}, {
					key: "viewToModel",
					value: function () {
						console.log("TODO")
					}
				}, {
					key: "prepareEvents",
					value: function (e, t) {
						if (e) {
							x.default.isFunction(e) && (e = e.call(this));
							var n, r = {};
							for (n in e) ! function () {
								var i = e[n];
								r[n] = x.default.isString(i) ? i : function (e) {
									var n = this.getItemByEv(e, !0);
									if (t) var r = i.call(this, n);
									else var r = i.call(this, e, n);
									return !1 !== r
								}
							}();
							return r
						}
					}
				}, {
					key: "getEvents",
					value: function () {
						var e = this.options,
							t = e.purist,
							n = e.events,
							r = e.ievents;
						n = this.prepareEvents(n, t), r = this.prepareEvents(r, !0);
						var i = this.getBaseEvents();
						return x.default.extend({}, i, n, r)
					}
				}, {
					key: "setSortMode",
					value: function (e) {
						this.options.sortMode = e, this.table.storePreference("sort-mode", e)
					}
				}, {
					key: "setViewType",
					value: function (e) {
						this.options.view = e, this.table.storePreference("view-type", e), this.table.render()
					}
				}, {
					key: "getColumnsVisibility",
					value: function () {
						var e = S.default.findByClass(this.rootElement, "visibility-check");
						return x.default.map(e, function (e) {
							return {
								name: S.default.attr(e, "name"),
								visible: e.checked
							}
						})
					}
				}, {
					key: "onColumnVisibilityChange",
					value: function () {
						var e = this.getColumnsVisibility();
						this.table.toggleColumns(e)
					}
				}, {
					key: "onViewChange",
					value: function (e) {
						if (!e) return !0;
						var t = e.target;
						this.setViewType(t.value)
					}
				}, {
					key: "onSortModeChange",
					value: function (e) {
						if (!e) return !0;
						var t = e.target;
						this.setSortMode(t.value)
					}
				}, {
					key: "onExportClick",
					value: function (e) {
						var t = e.target,
							n = t.dataset.format;
						n || (0, w.default)(29, "Missing format in export element's dataset."), this.table.exportTo(n)
					}
				}, {
					key: "getBaseEvents",
					value: function () {
						var e = {
							"click .pagination-bar-first-page": "goToFirst",
							"click .pagination-bar-last-page": "goToLast",
							"click .pagination-bar-prev-page": "goToPrev",
							"click .pagination-bar-next-page": "goToNext",
							"click .pagination-bar-refresh": "refresh",
							"change .pagination-bar-page-number": "changePage",
							"change .pagination-bar-results-select": "changeResultsNumber",
							"click .kt-advanced-filters": "toggleAdvancedFilters",
							"click .btn-clear-filters": "clearFilters",
							"click .king-table-head th.sortable": "sort",
							"keyup .search-field": "onSearchKeyUp",
							"paste .search-field, cut .search-field": "onSearchChange",
							"keyup .filters-region input[type='text']": "viewToModel",
							"keyup .filters-region textarea": "viewToModel",
							"change .filters-region input[type='checkbox']": "viewToModel",
							"change .filters-region input[type='radio']": "viewToModel",
							"change .filters-region select": "viewToModel",
							"change .visibility-check": "onColumnVisibilityChange",
							"click .export-btn": "onExportClick",
							"change [name='kt-view-type']": "onViewChange",
							"change [name='kt-sort-mode']": "onSortModeChange"
						};
						return x.default.each("text date datetime datetime-local email tel time search url week color month number".split(" "), function (t) {
							e["change .filters-region input[type='" + t + "']"] = "viewToModel"
						}), this.options.onItemClick && (e["click .kt-item"] = "onItemClick"), e
					}
				}, {
					key: "bindEvents",
					value: function () {
						var e = "__events__bound";
						return this[e] ? this : (this[e] = 1, this.delegateEvents().bindWindowEvents())
					}
				}, {
					key: "anyMenuIsOpen",
					value: function () {
						return !1
					}
				}, {
					key: "bindWindowEvents",
					value: function () {
						if ("undefined" != typeof window) {
							var e = this.unbindWindowEvents();
							S.default.on(document.body, "keydown.king-table", function (t) {
								if (S.default.anyInputFocused() || e.anyMenuIsOpen()) return !0;
								var n = t.keyCode;
								x.default.contains([37, 65], n) && e.goToPrev(), x.default.contains([39, 68], n) && e.goToNext()
							})
						}
						return this
					}
				}, {
					key: "unbindWindowEvents",
					value: function () {
						if ("undefined" != typeof window) {
							S.default.off(document.body, "keydown.king-table")
						}
						return this
					}
				}, {
					key: "delegateEvents",
					value: function () {
						var e = this,
							t = (e.table, e.options),
							n = e.table.element,
							r = e.getEvents();
						e.undelegateEvents();
						for (var i in r) {
							var a = r[i],
								s = a;
							if (s || (0, w.default)(27, "Invalid method definition"), x.default.isFunction(s) || (s = e[s]), !s && x.default.isFunction(t[a]) && (s = t[a]), !x.default.isFunction(s)) throw new Error("method not defined inside the model: " + r[i]);
							var o = i.match(/^(\S+)\s*(.*)$/),
								u = o[1],
								l = o[2];
							if (s = x.default.bind(s, e), u += ".delegate", "" === l) throw new Error("delegates without selector are not implemented");
							S.default.on(n, u, l, s)
						}
						return e.__events__bound = 0, e
					}
				}, {
					key: "undelegateEvents",
					value: function () {
						return S.default.off(this.table.element), this
					}
				}, {
					key: "dispose",
					value: function () {
						console.info("[*] RHTML DISPOSED!"), this.undelegateEvents().unbindWindowEvents(), S.default.remove(this.rootElement), S.default.removeClass(this.table.element, "king-table"), this.currentItems = this.rootElement = null, l(t.prototype.__proto__ || Object.getPrototypeOf(t.prototype), "dispose", this).call(this)
					}
				}, {
					key: "emptyView",
					value: function () {
						var e = this.getReg();
						return new c.VHtmlElement("div", {
							class: "king-table-empty"
						}, new c.VHtmlElement("span", 0, new c.VTextElement(e.noData)))
					}
				}, {
					key: "errorView",
					value: function (e) {
						return e || (e = this.getReg().errorFetchingData), new c.VHtmlFragment('<div class="king-table-error">\n      <span class="message">\n        <span>' + e + '</span>\n        <span class="oi" data-glyph="warning" aria-hidden="true"></span>\n      </span>\n    </div>')
					}
				}, {
					key: "loadingView",
					value: function () {
						var e = this.getReg();
						return new c.VHtmlElement("div", {
							class: "loading-info"
						}, [new c.VHtmlElement("span", {
							class: "loading-text"
						}, new c.VTextElement(e.loading)), new c.VHtmlElement("span", {
							class: "mini-loader"
						})])
					}
				}, {
					key: "singleLine",
					value: function () {
						throw new Error("make targeted updates")
					}
				}], [{
					key: "BaseHtmlBuilder",
					get: function () {
						return g.default
					}
				}, {
					key: "DomUtils",
					get: function () {
						return S.default
					}
				}]), t
			}(h.default);
		O.defaults = {
			view: "table",
			views: [{
				name: "table",
				resolver: !0
			}, {
				name: "gallery",
				resolver: _
			}],
			filtersView: null,
			filtersViewExpandable: !0,
			filtersViewOpen: !1,
			searchDelay: 50,
			sortMode: C.Simple,
			allowSortModes: !0,
			purist: !1,
			resultsPerPageSelect: [10, 30, 50, 100, 200],
			tools: null,
			prepTools: null,
			autoHighlightSearchProperties: !0
		}, n.default = O
	}, {
		"../../scripts/data/file": 11,
		"../../scripts/data/html": 12,
		"../../scripts/dom": 19,
		"../../scripts/menus/kingtable.menu": 25,
		"../../scripts/menus/kingtable.menu.html": 24,
		"../../scripts/raise": 26,
		"../../scripts/tables/kingtable.html.base.builder": 28,
		"../../scripts/tables/kingtable.html.builder": 29,
		"../../scripts/utils": 34
	}],
	33: [function (e, t, n) {
		"use strict";

		function r(e) {
			return e && e.__esModule ? e : {
				default: e
			}
		}

		function i(e, t) {
			if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
		}

		function a(e, t) {
			if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
			return !t || "object" != typeof t && "function" != typeof t ? e : t
		}

		function s(e, t) {
			if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
			e.prototype = Object.create(t && t.prototype, {
				constructor: {
					value: e,
					enumerable: !1,
					writable: !0,
					configurable: !0
				}
			}), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var o = function () {
				function e(e, t) {
					for (var n = 0; n < t.length; n++) {
						var r = t[n];
						r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
					}
				}
				return function (t, n, r) {
					return n && e(t.prototype, n), r && e(t, r), t
				}
			}(),
			u = e("../../scripts/tables/kingtable.builder"),
			l = r(u),
			f = e("../../scripts/literature/text-slider"),
			c = e("../../scripts/exceptions"),
			d = e("../../scripts/raise"),
			p = (r(d), e("../../scripts/utils")),
			h = r(p),
			m = e("../../scripts/components/string"),
			g = r(m),
			v = "\r\n",
			y = "*****************************************************************",
			b = function (e) {
				function t(e) {
					i(this, t);
					var n = a(this, (t.__proto__ || Object.getPrototypeOf(t)).call(this, e));
					return n.slider = new f.TextSlider("...."), n.setListeners(e), n
				}
				return s(t, e), o(t, [{
					key: "setListeners",
					value: function (e) {
						if (e) {
							var n = this;
							e.element && t.options.handleLoadingInfo && (n.listenTo(e, "fetching:data", function () {
								n.loadingHandler(e)
							}), n.listenTo(e, "fetched:data", function () {
								n.unsetLoadingHandler()
							}), n.listenTo(e, "fetch:fail", function () {
								n.unsetLoadingHandler().display(e, n.errorView())
							}), n.listenTo(e, "no-results", function () {
								n.unsetLoadingHandler().display(e, n.emptyView())
							}))
						}
					}
				}, {
					key: "singleLine",
					value: function (e, t) {
						return this.tabulate([[t]], [], h.default.extend({
							caption: e.options.caption
						}, e.pagination.totalPageCount > 0 ? e.pagination.data() : null))
					}
				}, {
					key: "errorView",
					value: function () {
						var e = this.table,
							t = this.getReg();
						return this.singleLine(e, t.errorFetchingData)
					}
				}, {
					key: "loadingHandler",
					value: function (e) {
						var n = this;
						n.unsetLoadingHandler();
						var r = this.slider,
							i = this.getReg(),
							a = i.loading + " ",
							s = e.hasData() ? t.options.loadInfoDelay : 0,
							o = e.element;
						n.showLoadingTimeout = setTimeout(function () {
							if (!e.loading) return n.unsetLoadingHandler();
							var t = n.singleLine(e, a + r.next());
							o && (o.innerHTML = t), n.loadingInterval = setInterval(function () {
								if (!e.loading) return n.unsetLoadingHandler();
								var t = n.singleLine(e, a + r.next());
								o.innerHTML = t
							}, 600)
						}, s)
					}
				}, {
					key: "unsetLoadingHandler",
					value: function () {
						return clearInterval(this.loadingInterval), clearTimeout(this.showLoadingTimeout), this.loadingInterval = this.showLoadingTimeout = null, this
					}
				}, {
					key: "dispose",
					value: function () {
						var e = table.element;
						e && (e.innerHTML = ""), this.stopListening(this.table), this.table = null, this.slider = null
					}
				}, {
					key: "display",
					value: function (e, t) {
						var n = e.element;
						if (n) {
							for (e.emit("empty:element", n); n.hasChildNodes();) n.removeChild(n.lastChild);
							if ("PRE" != n.tagName) {
								var r = document.createElement("pre");
								n.appendChild(r), n = r
							}
							n.innerHTML = t
						}
					}
				}, {
					key: "build",
					value: function (e) {
						e || (e = this.table);
						var t = e.getData({
							optimize: !0,
							format: !0
						});
						if (!t || !t.length) return this.display(e, this.emptyView());
						var n = t.shift(),
							r = this.tabulate(n, t, h.default.extend({
								caption: e.options.caption,
								dataAnchorTime: e.options.showAnchorTimestamp ? e.getFormattedAnchorTime() : null
							}, e.pagination.data()));
						this.display(e, r)
					}
				}, {
					key: "emptyView",
					value: function () {
						var e = this.getReg();
						return y + v + e.noData + v + y
					}
				}, {
					key: "paginationInfo",
					value: function (e) {
						var t = "",
							n = this.getReg(),
							r = e.page,
							i = e.totalPageCount,
							a = e.firstObjectNumber,
							s = e.lastObjectNumber,
							o = e.totalItemsCount,
							u = e.dataAnchorTime,
							l = h.default.isNumber;
						return l(r) && (t += n.page + " " + r, l(i) && (t += " " + n.of + " " + i), l(a) && l(s) && (t += " - " + n.results + g.default.format(" {0} - {1}", a, s), l(o) && (t += g.default.format(" " + n.of + " {0}", o)))), u && (t += " - " + n.anchorTime + " " + u), t
					}
				}, {
					key: "checkValue",
					value: function (e) {
						return e ? ("string" != typeof e && (e = e.toString()), e ? e.replace(/\r/g, "␍").replace(/\n/g, "␊") : "") : ""
					}
				}, {
					key: "tabulate",
					value: function (e, n, r) {
						var i = this;
						h.default.isArray(e) || (0, c.TypeException)("headers", "array"), n || (0, c.TypeException)("rows", "array"), h.default.any(n, function (e) {
							return !h.default.isArray(e)
						}) && (0, c.TypeException)("rows child", "array");
						var a = h.default.extend({}, t.options, r || {}),
							s = a.headersAlignment,
							o = a.rowsAlignment,
							u = a.padding,
							l = a.cornerChar,
							f = a.headerLineSeparator,
							d = a.cellVerticalLine,
							p = a.cellHorizontalLine,
							m = a.minCellWidth,
							y = a.headerCornerChar;
						u < 0 && (0, c.OutOfRangeException)("padding", 0);
						var b = this,
							w = e.length;
						w || (0, c.ArgumentException)("headers must contain at least one item"), h.default.any(n, function (e) {
							e.length
						}) && (0, c.ArgumentException)("each row must contain the same number of items");
						var k = "";
						h.default.reach(e, function (e) {
							return i.checkValue(e)
						}), h.default.reach(n, function (e) {
							return i.checkValue(e)
						});
						var x = g.default.ofLength(" ", u);
						u *= 2;
						var E, S = h.default.cols([e].concat(n)),
							P = h.default.map(S, function (e) {
								return Math.max(h.default.max(e, function (e) {
									return e.length
								}), m) + u
							}),
							T = a.caption,
							_ = 0;
						T && (_ = u + T.length + 2);
						var C = b.paginationInfo(a);
						if (C) {
							var O = u + C.length + 2;
							_ = Math.max(_, O)
						}
						if (_ > 0 && (E = h.default.sum(P) + P.length + 1, _ > E)) {
							var V = _ - E;
							P[P.length - 1] += V, E = _
						}
						T && (k += d + x + T + g.default.ofLength(" ", E - T.length - 3) + d + v), C && (k += d + x + C + g.default.ofLength(" ", E - C.length - 3) + d + v);
						var j = h.default.map(P, function (e) {
								return g.default.ofLength(f, e)
							}),
							M = h.default.map(P, function (e) {
								return g.default.ofLength(p, e)
							}),
							F = "";
						h.default.each(e, function (e, t) {
							F += y + j[t]
						}), F += y + v, (C || T) && (k = F + k), k += F, h.default.each(e, function (e, t) {
							k += d + b.align(x + e, P[t], s)
						}), k += d + v, k += F;
						var D, H = "";
						for (D = 0; D < w; D++) H += l + M[D];
						H += l;
						var I, N, A, L = n.length;
						for (D = 0; D < L; D++) {
							for (N = n[D], I = 0; I < w; I++) A = N[I], k += d + b.align(x + A, P[I], o);
							k += d + v, k += H + v
						}
						return k
					}
				}, {
					key: "align",
					value: function (e, t, n, r) {
						switch (r || (r = " "), n || (0, c.ArgumentNullException)("alignment"), n) {
							case "c":
							case "center":
								return g.default.center(e, t, r);
							case "l":
							case "left":
								return g.default.ljust(e, t, r);
							case "r":
							case "right":
								return g.default.rjust(e, t, r);
							default:
								(0, c.ArgumentException)("alignment: " + n)
						}
					}
				}], [{
					key: "options",
					get: function () {
						return {
							headerLineSeparator: "=",
							headerCornerChar: "=",
							cornerChar: "+",
							headersAlignment: "l",
							rowsAlignment: "l",
							padding: 1,
							cellVerticalLine: "|",
							cellHorizontalLine: "-",
							minCellWidth: 0,
							handleLoadingInfo: !0,
							loadInfoDelay: 500
						}
					}
				}]), t
			}(l.default);
		n.default = b
	}, {
		"../../scripts/components/string": 7,
		"../../scripts/exceptions": 20,
		"../../scripts/literature/text-slider": 23,
		"../../scripts/raise": 26,
		"../../scripts/tables/kingtable.builder": 27,
		"../../scripts/utils": 34
	}],
	34: [function (e, t, n) {
		"use strict";

		function r(e, t) {
			if ((!e || !e[F]) && p(e)) {
				var n, r = [];
				for (n in e) r.push(t(n, e[n]));
				return r
			}
			for (var r = [], i = 0, a = e[F]; i < a; i++) r.push(t(e[i]));
			return r
		}

		function i(e, t) {
			if (p(e)) {
				for (var n in e) t(e[n], n);
				return e
			}
			if (!e || !e[F]) return e;
			for (var r = 0, i = e[F]; r < i; r++) t(e[r], r)
		}

		function a(e, t) {
			for (var n = 0; n < t; n++) e(n)
		}

		function s(e) {
			return (void 0 === e ? "undefined" : C(e)) == V
		}

		function o(e) {
			return !isNaN(e) && (void 0 === e ? "undefined" : C(e)) == j
		}

		function u(e) {
			return (void 0 === e ? "undefined" : C(e)) == M
		}

		function l(e) {
			return (void 0 === e ? "undefined" : C(e)) == O
		}

		function f(e) {
			return e instanceof Array
		}

		function c(e) {
			return e instanceof Date
		}

		function d(e) {
			return e instanceof RegExp
		}

		function p(e) {
			return (void 0 === e ? "undefined" : C(e)) == O && null !== e && e.constructor == Object
		}

		function h(e) {
			if (!e) return !0;
			if (f(e)) return 0 == e.length;
			if (p(e)) {
				var t;
				for (t in e) return !1;
				return !0
			}
			if (s(e)) return "" === e;
			if (o(e)) return 0 === e;
			throw new Error("invalid argument")
		}

		function m(e, t) {
			return e && e.hasOwnProperty(t)
		}

		function g(e) {
			return e.toUpperCase()
		}

		function v(e) {
			return e.toLowerCase()
		}

		function y(e, t) {
			if (!t) return e ? e[0] : void 0;
			for (var n = 0, r = e[F]; n < r; n++)
				if (t(e[n])) return e[n]
		}

		function b(e) {
			return f(e) ? e : (void 0 === e ? "undefined" : C(e)) == O && e[F] ? r(e, function (e) {
				return e
			}) : Array.prototype.slice.call(arguments)
		}

		function w(e) {
			return f(e) ? [].concat.apply([], r(e, w)) : e
		}

		function k(e) {
			return D++, (e || "id") + D
		}

		function x() {
			D = -1
		}

		function E(e) {
			if (!e) return [];
			var t, n = [];
			for (t in e) n.push(t);
			return n
		}

		function S(e) {
			if (!e) return [];
			var t, n = [];
			for (t in e) n.push(e[t]);
			return n
		}

		function P(e, t) {
			if (!e) return e;
			t || (t = []);
			var n, r = {};
			for (n in e) - 1 == t.indexOf(n) && (r[n] = e[n]);
			return r
		}

		function T(e) {
			return void 0 === e
		}

		function _(e) {
			var t, n;
			if (null === e) return null;
			if (void 0 !== e) {
				if (l(e))
					if (f(e)) {
						n = [];
						for (var r = 0, i = e.length; r < i; r++) n[r] = _(e[r])
					} else {
						n = {};
						var a;
						for (t in e)
							if (null !== (a = e[t]) && void 0 !== a)
								if (l(a))
									if (c(a)) n[t] = new Date(a.getTime());
									else if (d(a)) n[t] = new RegExp(a.source, a.flags);
						else if (f(a)) {
							n[t] = [];
							for (var r = 0, i = a.length; r < i; r++) n[t][r] = _(a[r])
						} else n[t] = _(a);
						else n[t] = a;
						else n[t] = a
					}
				else n = e;
				return n
			}
		}
		Object.defineProperty(n, "__esModule", {
			value: !0
		});
		var C = "function" == typeof Symbol && "symbol" == typeof Symbol.iterator ? function (e) {
				return typeof e
			} : function (e) {
				return e && "function" == typeof Symbol && e.constructor === Symbol && e !== Symbol.prototype ? "symbol" : typeof e
			},
			O = (e("../scripts/exceptions"), "object"),
			V = "string",
			j = "number",
			M = "function",
			F = "length",
			D = -1;
		n.default = {
			extend: function () {
				var e = arguments;
				if (e[F]) {
					if (1 == e[F]) return e[0];
					for (var t, n, r = e[0], i = 1, a = e[F]; i < a; i++)
						if (t = e[i])
							for (n in t) r[n] = t[n];
					return r
				}
			},
			stringArgs: function (e) {
				if (!e || T(e.length)) throw new Error("expected array argument");
				if (!e.length) return [];
				if (1 === e.length) {
					var t = e[0];
					if (s(t) && t.indexOf(" ") > -1) return t.split(/\s+/g)
				}
				return e
			},
			uniqueId: k,
			resetSeed: x,
			flatten: w,
			each: i,
			exec: a,
			keys: E,
			values: S,
			minus: P,
			map: r,
			first: y,
			toArray: b,
			isArray: f,
			isDate: c,
			isString: s,
			isNumber: o,
			isObject: l,
			isPlainObject: p,
			isEmpty: h,
			isFunction: u,
			has: m,
			isNullOrEmptyString: function (e) {
				return null === e || void 0 === e || "" === e
			},
			lower: v,
			upper: g,
			clone: _,
			quacksLikePromise: function (e) {
				return !(!e || C(e.then) != M)
			},
			sum: function (e, t) {
				if (e) {
					var n, r = e[F];
					if (r) {
						for (var i = 0, r = e[F]; i < r; i++) {
							var a = t ? t(e[i]) : e[i];
							T(n) ? n = a : n += a
						}
						return n
					}
				}
			},
			max: function (e, t) {
				for (var n = -1 / 0, r = 0, i = e[F]; r < i; r++) {
					var a = t ? t(e[r]) : e[r];
					a > n && (n = a)
				}
				return n
			},
			min: function (e, t) {
				for (var n = 1 / 0, r = 0, i = e[F]; r < i; r++) {
					var a = t ? t(e[r]) : e[r];
					a < n && (n = a)
				}
				return n
			},
			withMax: function (e, t) {
				for (var n, r = 0, i = e[F]; r < i; r++)
					if (n) {
						var a = t(e[r]);
						a > t(n) && (n = e[r])
					} else n = e[r];
				return n
			},
			withMin: function (e, t) {
				for (var n, r = 0, i = e[F]; r < i; r++)
					if (n) {
						var a = t(e[r]);
						a < t(n) && (n = e[r])
					} else n = e[r];
				return n
			},
			indexOf: function (e, t) {
				return e.indexOf(t)
			},
			contains: function (e, t) {
				return e.indexOf(t) > -1
			},
			any: function (e, t) {
				if (p(e)) {
					var n;
					for (n in e)
						if (t(n, e[n])) return !0;
					return !1
				}
				for (var r = 0, i = e[F]; r < i; r++)
					if (t(e[r])) return !0;
				return !1
			},
			all: function (e, t) {
				if (p(e)) {
					var n;
					for (n in e)
						if (!t(n, e[n])) return !1;
					return !0
				}
				for (var r = 0, i = e[F]; r < i; r++)
					if (!t(e[r])) return !1;
				return !0
			},
			find: function (e, t) {
				if (!e) return null;
				if (f(e)) {
					if (!e || !e[F]) return;
					for (var n = 0, r = e[F]; n < r; n++)
						if (t(e[n])) return e[n]
				}
				if (p(e)) {
					var i;
					for (i in e)
						if (t(e[i], i)) return e[i]
				}
			},
			where: function (e, t) {
				if (!e || !e[F]) return [];
				for (var n = [], r = 0, i = e[F]; r < i; r++) t(e[r]) && n.push(e[r]);
				return n
			},
			removeItem: function (e, t) {
				for (var n = -1, r = 0, i = e[F]; r < i; r++)
					if (e[r] === t) {
						n = r;
						break
					}
				e.splice(n, 1)
			},
			reject: function (e, t) {
				if (!e || !e[F]) return [];
				for (var n = [], r = 0, i = e[F]; r < i; r++) t(e[r]) || n.push(e[r]);
				return n
			},
			pick: function (e, t, n) {
				var r = {};
				if (n)
					for (var i in e) - 1 == t.indexOf(i) && (r[i] = e[i]);
				else
					for (var a = 0, s = t[F]; a < s; a++) {
						var o = t[a];
						m(e, o) && (r[o] = e[o])
					}
				return r
			},
			require: function (e, t, n) {
				n || (n = "options");
				var r = "";
				if (e ? this.each(t, function (t) {
						m(e, t) || (r += "missing '" + t + "' in " + n)
					}) : r = "missing " + n, r) throw new Error(r)
			},
			wrap: function (e, t, n) {
				var r = this,
					i = arguments,
					a = function () {
						return t.apply(r, [e].concat(b(i)))
					};
				return a.bind(n || this), a
			},
			unwrap: function (e) {
				function t(t) {
					return e.apply(this, arguments)
				}
				return t.toString = function () {
					return e.toString()
				}, t
			}(function (e) {
				return u(e) ? unwrap(e()) : e
			}),
			defer: function (e) {
				setTimeout(e, 0)
			},
			atMost: function (e, t, n) {
				function r() {
					return e > 0 && (e--, i = t.apply(n || this, arguments)), i
				}
				var i;
				return r
			},
			isUnd: T,
			once: function (e, t) {
				return this.atMost(1, e, t)
			},
			partial: function (e) {
				var t = this,
					n = t.toArray(arguments);
				return n.shift(),
					function () {
						var r = t.toArray(arguments);
						return e.apply({}, n.concat(r))
					}
			},
			equal: function (e, t) {
				var n, r = !0;
				if (e === t) return r;
				if (e === n || t === n || null === e || null === t || e === r || t === r || !1 === e || !1 === t || "" === e || "" === t) return !1;
				if (f(e)) {
					if (f(t) && e[F] == t[F]) {
						var i, a = e[F];
						for (i = 0; i < a; i++)
							if (!this.equal(e[i], t[i])) return !1;
						return r
					}
					return !1
				}
				if (o(e) || s(e)) return e == t;
				if (null === e && null === t) return r;
				if (e === n && t === n) return r;
				var u, l = 0,
					c = 0;
				for (u in e)
					if (e[u] !== n && (l += 1), !this.equal(e[u], t[u])) return !1;
				for (u in t) t[u] !== n && (c += 1);
				return l == c
			},
			cols: function (e) {
				if (!e || !e.length) return [];
				var t, n, r = this.max(e, function (e) {
						return e.length
					}),
					i = [],
					a = e.length;
				for (n = 0; n < r; n++) {
					var s = [];
					for (t = 0; t < a; t++) s.push(e[t][n]);
					i.push(s)
				}
				return i
			},
			sortNums: function (e) {
				return e.sort(function (e, t) {
					return e > t ? 1 : e < t ? -1 : 0
				})
			},
			debounce: function (e, t, n) {
				function r() {
					i && clearTimeout(i);
					var r = arguments.length ? b(arguments) : void 0;
					i = setTimeout(function () {
						i = null, e.apply(n, r)
					}, t)
				}
				var i;
				return r
			},
			reach: function (e, t) {
				if (!f(e)) throw new Error("expected array");
				for (var n, r = 0, i = e.length; r < i; r++) n = e[r], f(n) ? this.reach(n, t) : e[r] = t(n);
				return e
			},
			quacks: function (e, t) {
				if (!e) return !1;
				if (!t) throw "missing methods list";
				s(t) && (t = b(arguments).slice(1, arguments.length));
				for (var n = 0, r = t.length; n < r; n++)
					if (!u(e[t[n]])) return !1;
				return !0
			},
			format: function (e, t) {
				return e.replace(/\{\{(.+?)\}\}/g, function (e, n) {
					return t.hasOwnProperty(n) ? t[n] : e
				})
			},
			bind: function (e, t) {
				return e.bind(t)
			},
			ifcall: function (e, t, n) {
				if (e) {
					if (!n) return void e.call(t);
					switch (n.length) {
						case 0:
							return void e.call(t);
						case 1:
							return void e.call(t, n[0]);
						case 2:
							return void e.call(t, n[0], n[1]);
						case 3:
							return void e.call(t, n[0], n[1], n[2]);
						default:
							e.apply(t, n)
					}
				}
			}
		}
	}, {
		"../scripts/exceptions": 20
	}]
}, {}, [30]);
