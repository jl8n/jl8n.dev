<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/stores";

    let drawerIsOpen = false;

    function toggleDrawer() {
        drawerIsOpen = !drawerIsOpen;
    }

    function closeDrawer() {
        drawerIsOpen = !drawerIsOpen;
    }

    onMount(() => {
        document.addEventListener("keydown", (event) => {
            if (event.key === "Escape") {
                toggleDrawer();
            }
        });
    });
</script>

<button
    on:click={toggleDrawer}
    type="button"
    class="menu-button material-symbols-outlined"
>
    menu
</button>

<div class="drawer" aria-expanded={drawerIsOpen}>
    <button
        type="button"
        class="drawer-close material-symbols-outlined"
        on:click={toggleDrawer}>close</button
    >

    <div class="drawer-wrapper">
        <div class="drawer-content">
            <nav>
                <a
                    href="/"
                    on:click={closeDrawer}
                    aria-current={$page.url.pathname === "/"
                        ? "page"
                        : undefined}>Home</a
                >
                <a
                    href="/projects"
                    on:click={closeDrawer}
                    aria-current={$page.url.pathname === "/projects"
                        ? "page"
                        : undefined}>Projects</a
                >
                <a
                    href="/listens"
                    on:click={closeDrawer}
                    aria-current={$page.url.pathname === "/listens"
                        ? "page"
                        : undefined}>Listens</a
                >
            </nav>
        </div>
    </div>
</div>

<style lang="scss">
    button {
        font-size: 40px;
        color: white;

        background-color: transparent;
        border: none;
        cursor: pointer;
    }

    .drawer {
        position: fixed;
        top: 0;
        left: -100vw;
        width: 100vw;
        height: 100vh;
        background-color: rgb(51, 47, 47);
        z-index: 100;
        opacity: 0;
        visibility: hidden;
        pointer-events: none;
        transition: opacity 0.2s ease-in-out, left 0.2s ease-in-out,
            visibility 0.2s ease-in-out;

        &[aria-expanded="true"] {
            opacity: 1;
            visibility: visible;
            pointer-events: all;
            left: 0;
        }
    }

    .drawer-wrapper {
        position: absolute;
        color: black;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: 1;
    }

    nav {
        display: flex;
        flex-direction: column;
        width: 100%;
        font-size: 4em;
        padding-top: 1em;
        padding-left: 15px;
    }

    nav > a {
        text-decoration: none;

        &:hover {
            color: aqua;
        }
    }

    .drawer-close {
        position: absolute;
        top: 0px;
        right: 10px;
        color: black;
        z-index: 2;
        padding: 10px;
    }
</style>
