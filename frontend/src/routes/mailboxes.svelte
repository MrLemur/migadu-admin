<script>
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { activeDomain } from '../components/stores.js';
	import 'carbon-components-svelte/css/g80.css';
	import { slide } from 'svelte/transition';
	import Modal from '../components/ModalForm.svelte';

	let mailboxes = [];
	let identities = [];
	let loading = true;
	let activeMailbox;

	const MailboxModal = {
		isShown: false,
		data: {},
		Modal
	};

	const IdentityModal = {
		isShown: false,
		data: {},
		Modal
	};

	const toggleIdentities = (mailbox) => {
		mailboxes = mailboxes.map((m) => {
			if (m.address === mailbox.address) {
				m.active = !m.active;
				activeMailbox = mailbox;
			}
			return m;
		});
	};

	const getMailboxes = async () => {
		console.log(loading);
		loading = true;
		console.log(loading);
		mailboxes = [];
		mailboxes = mailboxes;
		const response = await fetch(`http://localhost:5000/api/${$activeDomain}/mailboxes`);
		mailboxes = await response.json();
		mailboxes.forEach((mailbox) => {
			mailbox.active = false;
			identities.push(mailbox.identities);
		});
		console.log(loading);
		loading = false;
		console.log(loading);
	};

	const newMailbox = async (localPart, displayName, invitationEmail, password) => {
		const response = await fetch(`http://localhost:5000/api/${domain}/mailboxes`, {
			method: 'POST',
			body: JSON.stringify({
				localPart,
				displayName,
				invitationEmail,
				password
			})
		});
		getMailboxes();
	};

	const newIdentity = async (mailbox, localPart, displayName) => {
		const response = await fetch(
			`http://localhost:5000/api/${$activeDomain}/identities/${mailbox}`,
			{
				method: 'POST',
				body: JSON.stringify({
					mailbox,
					localPart,
					displayName
				})
			}
		);
		getMailboxes();
	};

	const deleteMailbox = async (id) => {
		const response = await fetch(`/api/mailboxes/${id}`, {
			method: 'DELETE'
		});
		getMailboxes();
	};
	const deleteIdentity = async (mailbox, localPart) => {
		const response = await fetch(`/api/identities/${mailbox}/${localPart}`, {
			method: 'DELETE'
		});
		getMailboxes();
	};

	onMount(async () => {
		getMailboxes();
	});

	$: getMailboxes($activeDomain);
</script>

<div>
	<h1 color="white">{loading}</h1>
	{#if MailboxModal.isShown === true}
		<MailboxModal.Modal
			title={MailboxModal.data.title}
			formFields={MailboxModal.data.fields}
			on:formSubmitted={async (event) => {
				MailboxModal.isShown = true;
				const data = event.detail;
				await newMailbox(data.localPart, data.displayName, data.invitationEmail, data.password);
				MailboxModal.isShown = false;
				getMailboxes();
			}}
			on:formCancelled={() => {
				MailboxModal.isShown = false;
			}}
		/>
	{/if}
	{#if IdentityModal.isShown === true}
		<IdentityModal.Modal
			title={IdentityModal.data.title}
			formFields={IdentityModal.data.fields}
			on:formSubmitted={async (event) => {
				IdentityModal.isShown = true;
				const data = event.detail;
				const mailbox = activeMailbox;
				await newIdentity(mailbox.local_part, data.localPart, data.displayName);
				IdentityModal.isShown = false;
				getMailboxes();
			}}
			on:formCancelled={() => {
				IdentityModal.isShown = false;
			}}
		/>
	{/if}

	{#if !loading}
		<h1>HELLO</h1>
		<table>
			<tr>
				<td colspan="3" align="center"
					><button
						on:click={async () => {
							MailboxModal.data.title = `New Mailbox`;
							MailboxModal.data.fields = [
								{ label: 'Local Part', name: 'localPart', type: 'text' },
								{ label: 'Display Name', name: 'displayName', type: 'text' },
								{ label: 'Invitation Email', name: 'invitationEmail', type: 'text' },
								{ label: 'Password (optional)', name: 'password', type: 'password' }
							];
							MailboxModal.isShown = true;
						}}
						style="margin: auto;">New Mailbox</button
					>
				</td>
			</tr>
			<tr>
				<th>Email Address</th>
				<th>Display Name</th>
				<th>Actions</th>
			</tr>
			{#each mailboxes as mailbox (mailbox.address)}
				<tr>
					<td>{mailbox.address}</td>
					<td>{mailbox.name}</td>
					<td
						><button on:click={() => editMailbox(mailbox)}>Edit</button>
						<button on:click={() => toggleIdentities(mailbox)}>Identities</button>
						<button on:click={() => deleteMailbox(mailbox.local_part)}>Delete</button></td
					>
				</tr>
				{#if mailbox.active}
					<tr>
						<td colspan="3">
							<div
								transition:slide={{ duration: 300 }}
								style="background-color: teal; width: 100%;"
							>
								<table style=" margin: auto;">
									<tr>
										<td colspan="3" align="center"
											><button
												on:click={() => {
													IdentityModal.data.title = `New Identity`;
													IdentityModal.data.fields = [
														{
															label: 'Local Part',
															name: 'localPart',
															type: 'text',
															placeholder: `@${mailbox.domain_name}`
														},
														{ label: 'Display Name', name: 'displayName', type: 'text' }
													];
													IdentityModal.isShown = true;
												}}
												style="margin: auto;">New Identity</button
											>
										</td>
									</tr>
									{#if mailbox.identities}
										{#each mailbox.identities as identity}
											<tr>
												<td>{identity.address}</td>
												<td>{identity.name}</td>
												<td
													><button>Edit</button>
													<button
														on:click={() => deleteIdentity(mailbox.local_part, identity.local_part)}
														>Delete</button
													></td
												>
											</tr>
										{/each}
									{/if}
								</table>
							</div>
						</td>
					</tr>
				{/if}
			{:else}
				{console.log("Didn't work")}
				<p>Nothing</p>
			{/each}
		</table>
	{:else}
		<h3>Loading...</h3>
	{/if}
</div>

<style>
	table {
		padding: 10px;
	}

	td {
		padding: 5px;
	}

	.identity {
		background-color: aqua;
	}
</style>
