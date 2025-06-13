<div class="content">
  <div class="card-wrapper">
    <Card footer={false} fill={false}>
      <span slot="title">
        Admin Serverübersicht
      </span>

      <div slot="body" style="width: 100%">
        <div id="guild-container">
          {#each paginatedGuilds as guild (guild.id)}
            <Guild guild={guild}/>
          {/each}
        </div>

        <div class="flex-container" id="refresh-container" bind:this={refreshContainer}>
          <Button icon="fas fa-sync" on:click={refreshGuilds}>
            Liste Aktuallisieren
          </Button>
        </div>

        <div class="pagination-controls" bind:this={paginationControls}>
            <Button icon="fas fa-arrow-left" on:click={prevPage} disabled={currentPage === 1}>
                Zurück
            </Button>
            
            {#each Array(totalPages) as _, i}
                <Button
                    on:click={() => goToPage(i + 1)}
                    type="button"
                    active={currentPage === i + 1}
                    noShadow>
                        {i + 1}
                </Button>
            {/each}
            <Button icon="fas fa-arrow-right" on:click={nextPage} disabled={currentPage === totalPages}>
                Weiter
            </Button>
        </div>
      </div>
    </Card>
  </div>
</div>

<script>
    import axios from 'axios';
    import { onMount } from 'svelte';
    import {fade} from 'svelte/transition';
    import {notifyError, withLoadingScreen} from '../js/util'
    import {setDefaultHeaders} from '../includes/Auth.svelte'
    import {API_URL} from "../js/constants.js";
    import Guild from '../components/Guild.svelte'
    import Card from '../components/Card.svelte'
    import InviteBadge from '../components/InviteBadge.svelte'
    import Button from '../components/Button.svelte'
    import {loadingScreen, permissionLevelCache} from "../js/stores";

    setDefaultHeaders();

    let refreshContainer;
    let paginationControls;

    let guildsall = window.localStorage.getItem('guildsall') ? JSON.parse(window.localStorage.getItem('guildsall')) : [];
    let currentPage = 1;
    let itemsPerPage = 15; 
    
    $: totalPages = Math.ceil(guildsall.length / itemsPerPage);
    
    $: paginatedGuilds = guildsall.slice(
      (currentPage - 1) * itemsPerPage,
      currentPage * itemsPerPage
    );

    function nextPage() {
      if (currentPage < totalPages) currentPage++;
    }
    
    function prevPage() {
      if (currentPage > 1) currentPage--;
    }
    
    function goToPage(page) {
      if (page >= 1 && page <= totalPages) currentPage = page;
    }

    async function refreshGuilds() {
        await withLoadingScreen(async () => {
            const res = await axios.post(`${API_URL}/user/guilds/reloadall`);
            if (res.status !== 200) {
                notifyError(res.data.error);
                return;
            }

            if (!res.data.success && res.data['reauthenticate_required'] === true) {
                window.location.href = "/login";
                return;
            }

            guildsall = res.data.guilds;
            window.localStorage.setItem('guildsall', JSON.stringify(guildsall));
        });
    }

    function recalcItemsPerPage() {
        if (!guildContainer) return;

        const containerHeight = guildContainer.clientHeight;

        const badgeHeight = 110;

        const cardsPerRow = window.innerWidth > 950 ? 3 : 1;

        const reservedHeight =
            (refreshContainer?.offsetHeight || 0) +
            (paginationControls?.offsetHeight || 0);

        const usableHeight = containerHeight - reservedHeight;
        const rows = Math.floor(usableHeight / badgeHeight) || 1;

        itemsPerPage = cardsPerRow * rows;
        currentPage = 1;
    }

    onMount(() => {
        recalcItemsPerPage();

        window.addEventListener('resize', () => {
            recalcItemsPerPage();
        });
    });

    loadingScreen.set(false);
</script>

<style>
    .content {
        display: flex;
        height: 100%;
        width: 100%;
        justify-content: center;
    }

    .card-wrapper {
        display: block;
        width: 75%;
        height: 90%;
        margin-top: 5%;
        margin-bottom: 5%;
    }

    #guild-container {
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        justify-content: space-evenly;
    }

    #refresh-container {
        display: flex;
        justify-content: center;

        margin: 10px 0;
        color: white;
    }

    @media (max-width: 576px) {
        .card-wrapper {
            width: 100%;
        }
    }
    .pagination-controls {
        margin-top: 1rem;
        display: flex;
        gap: 0.5rem;
        justify-content: center;
    }

    .pagination-controls button.active {
        font-weight: bold;
        background: #ddd;
    }

    .pagination-controls button:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }
</style>
