(function () {
    const AVATAR_EMOJIS = [
        "😀", "😎", "🤩", "🥳", "😺", "🐶", "🐱", "🦊", "🐻", "🐼", "🐨", "🦁",
        "🐯", "🐸", "🐵", "🦄", "🐲", "🌟", "🔥", "💎", "🎮", "🎯", "👾", "🤖"
    ];

    let authToken = localStorage.getItem("game_token") || null;
    let currentUser = null;
    let mountId = "userArea";
    let gameKey = "";
    let authMode = "login";
    let onAuthSuccess = null;
    let selectedEmoji = "";

    function esc(text) {
        const div = document.createElement("div");
        div.textContent = text;
        return div.innerHTML;
    }

    function toast(msg, type) {
        const el = document.createElement("div");
        el.className = "toast toast-" + (type || "success");
        el.textContent = msg;
        document.body.appendChild(el);
        requestAnimationFrame(() => el.classList.add("show"));
        setTimeout(() => {
            el.classList.remove("show");
            setTimeout(() => el.remove(), 200);
        }, 2200);
    }

    function getAvatarDisplay(user) {
        if (!user) return "";
        if (user.avatar && user.avatar.startsWith("/")) {
            return '<img src="' + user.avatar + '" alt="">';
        }
        if (user.avatar && user.avatar.length <= 4) {
            return user.avatar;
        }
        return user.username ? user.username.charAt(0).toUpperCase() : "?";
    }

    async function api(path, opts) {
        const headers = {};
        if (!(opts && opts.body instanceof FormData)) {
            headers["Content-Type"] = "application/json";
        }
        if (authToken) {
            headers["Authorization"] = "Bearer " + authToken;
        }
        const response = await fetch(path, {
            ...opts,
            headers: { ...headers, ...(opts && opts.headers ? opts.headers : {}) }
        });
        return response.json();
    }

    function updateUserArea() {
        const mount = document.getElementById(mountId);
        if (!mount) return;

        if (!currentUser) {
            mount.innerHTML = '<button class="hdr-btn" id="btnQuickLogin">👤 登录</button>';
            const button = document.getElementById("btnQuickLogin");
            if (button) {
                button.addEventListener("click", () => showAuth("login"));
            }
            return;
        }

        mount.innerHTML =
            '<div class="avatar-wrapper">' +
                '<div class="avatar-btn" id="avatarBtn">' + getAvatarDisplay(currentUser) + '</div>' +
                '<div class="avatar-dropdown" id="avatarDD">' +
                    '<div class="dd-header">' +
                        '<div class="dd-name">' + esc(currentUser.username || "") + '</div>' +
                        '<div class="dd-phone">' + (currentUser.phone || "") + '</div>' +
                    '</div>' +
                    '<div class="dd-item" id="ddProfile">个人资料</div>' +
                    '<div class="dd-item" id="ddAvatar">修改头像</div>' +
                    '<div class="dd-divider"></div>' +
                    '<div class="dd-item danger" id="ddLogout">退出登录</div>' +
                '</div>' +
            '</div>';

        document.getElementById("avatarBtn").addEventListener("click", function (event) {
            event.stopPropagation();
            document.getElementById("avatarDD").classList.toggle("show");
        });
        document.getElementById("ddProfile").addEventListener("click", openProfile);
        document.getElementById("ddAvatar").addEventListener("click", openAvatarPicker);
        document.getElementById("ddLogout").addEventListener("click", logout);
    }

    function closeDropdown() {
        const dd = document.getElementById("avatarDD");
        if (dd) dd.classList.remove("show");
    }

    function ensureModals() {
        if (document.getElementById("authOverlay")) return;
        const wrapper = document.createElement("div");
        wrapper.innerHTML =
            '<div class="overlay" id="authOverlay"><div class="overlay-box">' +
                '<button class="close-x" id="closeAuthBtn">&times;</button>' +
                '<h2 id="authTitle">登录</h2>' +
                '<p id="authSubtitle"></p>' +
                '<div class="form-group"><input class="form-input" type="tel" id="phoneInput" placeholder="手机号" maxlength="11"><div class="form-err" id="phoneErr"></div></div>' +
                '<div class="form-group"><input class="form-input" type="password" id="pwdInput" placeholder="密码（至少 6 位）"><div class="form-err" id="pwdErr"></div></div>' +
                '<button class="btn btn-primary" style="width:100%" id="authSubmitBtn">登录</button>' +
                '<div class="auth-switch"><span id="authSwitchText">没有账号？</span> <a id="authSwitchLink">去注册</a></div>' +
            '</div></div>' +
            '<div class="overlay" id="profileOverlay"><div class="overlay-box">' +
                '<button class="close-x" id="closeProfileBtn">&times;</button>' +
                '<h2>个人资料</h2>' +
                '<div class="profile-avatar-lg" id="profileAvatar"></div>' +
                '<p style="font-size:12px;color:rgba(255,255,255,0.5)">点击头像可更换</p>' +
                '<div class="form-group"><input class="form-input" type="text" id="profileNameInput" maxlength="20" placeholder="用户名"></div>' +
                '<div class="form-group"><input class="form-input" type="text" id="profilePhoneInput" disabled style="opacity:0.55"></div>' +
                '<button class="btn btn-primary" style="width:100%" id="saveProfileBtn">保存</button>' +
            '</div></div>' +
            '<div class="overlay" id="avatarOverlay"><div class="overlay-box">' +
                '<button class="close-x" id="closeAvatarBtn">&times;</button>' +
                '<h2>选择头像</h2>' +
                '<div class="emoji-grid" id="emojiGrid"></div>' +
                '<div class="upload-zone" id="uploadZone">点击上传头像（最大 2MB）</div>' +
                '<input type="file" id="avatarFileInput" accept="image/*" style="display:none">' +
                '<button class="btn btn-primary" style="width:100%;margin-top:12px" id="confirmEmojiBtn">确认选择</button>' +
            '</div></div>';
        document.body.appendChild(wrapper);

        document.getElementById("closeAuthBtn").addEventListener("click", closeAuth);
        document.getElementById("authSwitchLink").addEventListener("click", toggleAuthMode);
        document.getElementById("authSubmitBtn").addEventListener("click", submitAuth);
        document.getElementById("pwdInput").addEventListener("keydown", function (e) {
            if (e.key === "Enter") submitAuth();
        });
        document.getElementById("closeProfileBtn").addEventListener("click", function () {
            document.getElementById("profileOverlay").classList.remove("active");
        });
        document.getElementById("profileAvatar").addEventListener("click", openAvatarPicker);
        document.getElementById("saveProfileBtn").addEventListener("click", saveProfile);
        document.getElementById("closeAvatarBtn").addEventListener("click", function () {
            document.getElementById("avatarOverlay").classList.remove("active");
        });
        document.getElementById("uploadZone").addEventListener("click", function () {
            document.getElementById("avatarFileInput").click();
        });
        document.getElementById("avatarFileInput").addEventListener("change", uploadAvatarFile);
        document.getElementById("confirmEmojiBtn").addEventListener("click", confirmEmojiAvatar);
    }

    function showErr(id, msg) {
        const el = document.getElementById(id);
        el.textContent = msg;
        el.style.display = "block";
    }
    function hideErr(id) {
        document.getElementById(id).style.display = "none";
    }

    function syncAuthModalText() {
        document.getElementById("authTitle").textContent = authMode === "login" ? "登录" : "注册";
        document.getElementById("authSubtitle").textContent = authMode === "login" ? "使用手机号登录" : "注册后自动登录";
        document.getElementById("authSubmitBtn").textContent = authMode === "login" ? "登录" : "注册";
        document.getElementById("authSwitchText").textContent = authMode === "login" ? "没有账号？" : "已有账号？";
        document.getElementById("authSwitchLink").textContent = authMode === "login" ? "去注册" : "去登录";
    }

    function showAuth(mode, successCb) {
        authMode = mode || "login";
        onAuthSuccess = typeof successCb === "function" ? successCb : null;
        syncAuthModalText();
        document.getElementById("phoneInput").value = "";
        document.getElementById("pwdInput").value = "";
        hideErr("phoneErr");
        hideErr("pwdErr");
        document.getElementById("authOverlay").classList.add("active");
    }

    function closeAuth() {
        document.getElementById("authOverlay").classList.remove("active");
    }

    function toggleAuthMode() {
        showAuth(authMode === "login" ? "register" : "login", onAuthSuccess);
    }

    async function submitAuth() {
        const phone = document.getElementById("phoneInput").value.trim();
        const pwd = document.getElementById("pwdInput").value;
        hideErr("phoneErr");
        hideErr("pwdErr");

        if (!/^1\d{10}$/.test(phone)) { showErr("phoneErr", "请输入 11 位手机号"); return; }
        if (pwd.length < 6) { showErr("pwdErr", "密码至少 6 位"); return; }

        const button = document.getElementById("authSubmitBtn");
        button.disabled = true;
        button.textContent = "请稍候...";

        try {
            const isLogin = authMode === "login";
            const path = isLogin ? "/api/auth/login" : "/api/auth/register";
            const body = isLogin
                ? { phone: phone, password: pwd }
                : { phone: phone, password: pwd, source: gameKey };
            const data = await api(path, { method: "POST", body: JSON.stringify(body) });
            if (data.code !== 0) {
                toast(data.msg || "认证失败", "error");
                return;
            }
            authToken = data.data.token;
            currentUser = data.data.user;
            localStorage.setItem("game_token", authToken);
            updateUserArea();
            closeAuth();
            toast(isLogin ? "登录成功" : "注册成功", "success");
            if (onAuthSuccess) {
                const cb = onAuthSuccess;
                onAuthSuccess = null;
                cb(currentUser);
            }
        } catch (error) {
            toast("网络错误", "error");
        } finally {
            button.disabled = false;
            button.textContent = authMode === "login" ? "登录" : "注册";
        }
    }

    async function saveProfile() {
        const username = document.getElementById("profileNameInput").value.trim();
        if (!username || username.length > 20) { toast("用户名长度为 1-20", "error"); return; }
        try {
            const data = await api("/api/user/profile", { method: "PUT", body: JSON.stringify({ username: username }) });
            if (data.code !== 0) { toast(data.msg || "更新失败", "error"); return; }
            currentUser = data.data.user;
            updateUserArea();
            document.getElementById("profileOverlay").classList.remove("active");
            toast("资料已更新", "success");
        } catch (error) { toast("网络错误", "error"); }
    }

    function openProfile() {
        closeDropdown();
        if (!currentUser) return;
        document.getElementById("profileAvatar").innerHTML = getAvatarDisplay(currentUser);
        document.getElementById("profileNameInput").value = currentUser.username || "";
        document.getElementById("profilePhoneInput").value = currentUser.phone || "";
        document.getElementById("profileOverlay").classList.add("active");
    }

    function openAvatarPicker() {
        closeDropdown();
        selectedEmoji = "";
        document.getElementById("profileOverlay").classList.remove("active");
        const grid = document.getElementById("emojiGrid");
        grid.innerHTML = AVATAR_EMOJIS.map(function (emoji) {
            return '<div class="emoji-opt" data-emoji="' + emoji + '">' + emoji + '</div>';
        }).join("");
        grid.querySelectorAll(".emoji-opt").forEach(function (node) {
            node.addEventListener("click", function () {
                grid.querySelectorAll(".emoji-opt").forEach(function (n) { n.classList.remove("selected"); });
                node.classList.add("selected");
                selectedEmoji = node.getAttribute("data-emoji");
            });
        });
        document.getElementById("avatarOverlay").classList.add("active");
    }

    async function confirmEmojiAvatar() {
        if (!selectedEmoji) { toast("请先选择一个头像", "error"); return; }
        try {
            const data = await api("/api/user/avatar", { method: "PUT", body: JSON.stringify({ avatar: selectedEmoji }) });
            if (data.code !== 0) { toast(data.msg || "更新失败", "error"); return; }
            currentUser = data.data.user;
            updateUserArea();
            document.getElementById("avatarOverlay").classList.remove("active");
            toast("头像已更新", "success");
        } catch (error) { toast("网络错误", "error"); }
    }

    async function uploadAvatarFile(event) {
        const file = event.target.files[0];
        if (!file) return;
        if (file.size > 2 * 1024 * 1024) { toast("图片不能超过 2MB", "error"); event.target.value = ""; return; }
        const form = new FormData();
        form.append("file", file);
        try {
            const data = await api("/api/user/avatar", { method: "POST", body: form });
            if (data.code !== 0) { toast(data.msg || "上传失败", "error"); return; }
            currentUser = data.data.user;
            updateUserArea();
            document.getElementById("avatarOverlay").classList.remove("active");
            toast("头像已更新", "success");
        } catch (error) { toast("上传失败", "error"); }
        finally { event.target.value = ""; }
    }

    function logout() {
        authToken = null;
        currentUser = null;
        localStorage.removeItem("game_token");
        updateUserArea();
        closeDropdown();
        toast("已退出登录", "success");
    }

    async function checkAuth() {
        if (!authToken) { updateUserArea(); return; }
        try {
            const data = await api("/api/user");
            if (data.code === 0) { currentUser = data.data.user; updateUserArea(); return; }
        } catch (error) {}
        authToken = null;
        currentUser = null;
        localStorage.removeItem("game_token");
        updateUserArea();
    }

    function requireLogin(successCb) {
        if (currentUser) {
            if (typeof successCb === "function") successCb(currentUser);
            return true;
        }
        showAuth("login", successCb);
        return false;
    }

    function init(options) {
        mountId = (options && options.mountId) ? options.mountId : "userArea";
        gameKey = (options && options.gameKey) ? options.gameKey : "";
        ensureModals();
        updateUserArea();
        checkAuth();
        document.addEventListener("click", closeDropdown);
    }

    window.GamePortalAuth = {
        init: init,
        api: api,
        showAuth: showAuth,
        requireLogin: requireLogin,
        getCurrentUser: function () { return currentUser; },
        toast: toast
    };
})();
