{#if data}
  <ConfirmationModal icon="fas fa-floppy-disk" on:confirm={() => dispatch("confirm", data)} on:cancel={() => dispatch("cancel", {})}>
    <span slot="title">Tag Bearbeiter</span>
    <div slot="body" class="body-wrapper">
      <div class="row">
        <Input col4 label="Tag ID" placeholder="docs" bind:value={data.id}
               tooltipText='Wenn der Befehl "/tag docs" ist, dann ist die ID "docs"'/>
        <Checkbox col2 label="Eigener Command Alias erstellen" bind:value={data.use_guild_command}
                  disabled={!isPremium} premiumBadge={true} />
      </div>

      <div class="row">
        <Textarea col1 label="Nachricht Inhalt" bind:value={data.content} placeholder="Nachricht Inhalt, auserhalb des embeds"/>
      </div>

      <div class="row">
      </div>

      <div class="col">
        <div class="inline">
          <Toggle inline label="Embed benutzen" bind:value={data.use_embed}/>
          <hr/>
        </div>

        {#if data.use_embed}
          <EmbedForm footerPremiumOnly={false} bind:data={data.embed}/>
        {/if}
      </div>
    </div>

    <span slot="confirm">Speichern</span>
  </ConfirmationModal>
{/if}

<svelte:window on:keydown={handleKeydown}/>

<script>
    import ConfirmationModal from "../ConfirmationModal.svelte";
    import Input from "../form/Input.svelte";
    import Checkbox from "../form/Checkbox.svelte";
    import Textarea from "../form/Textarea.svelte";
    import Toggle from "../form/Toggle.svelte";
    import EmbedForm from "../EmbedForm.svelte";
    import {createEventDispatcher, onMount} from "svelte";

    const dispatch = createEventDispatcher();

    export let data;
    export let isPremium;

    function handleKeydown(e) {
        if (e.key === "Escape") {
            dispatch("cancel", {});
        }
    }

    onMount(() => {
        if (data === undefined) {
            data = {
                use_embed: false,
            };
        }
    })
</script>

<style>
    .body-wrapper {
        display: flex;
        flex-direction: column;
        width: 100%;
    }

    .row {
        display: flex;
        flex-direction: row;
        gap: 2%;
    }

    .col {
        display: flex;
        flex-direction: column;
        row-gap: 2vh;
    }

    .inline {
        display: flex;
        flex-direction: row;
        align-items: center;
        width: 100%;
        gap: 10px;
    }

    hr {
        border-top: 1px solid #777;
        border-bottom: 0;
        border-left: 0;
        border-right: 0;
        width: 100%;
        flex: 1;
    }
</style>
